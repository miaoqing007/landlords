package wsconnection

import (
	"errors"
	"github.com/gorilla/websocket"
	"landlords/manager"
	"landlords/registry"
	"log"
	"sync"
)

// ws 的所有连接
// 用于广播
var WsConnAll map[int64]*WsConnection

// 客户端连接
type WsConnection struct {
	WsSocket *websocket.Conn // 底层websocket
	InChan   chan []byte     // 读队列
	OutChan  chan []byte     // 写队列

	mutex     sync.Mutex // 避免重复关闭管道,加锁处理
	isClosed  bool
	CloseChan chan byte // 关闭通知
	id        int64
	*manager.Player
}

func NewWsConnection(wsSocket *websocket.Conn, maxConnId int64) *WsConnection {
	wsConn := &WsConnection{
		WsSocket:  wsSocket,
		InChan:    make(chan []byte, 1000),
		OutChan:   make(chan []byte, 1000),
		CloseChan: make(chan byte),
		isClosed:  false,
		id:        maxConnId,
	}
	WsConnAll[maxConnId] = wsConn
	log.Println("当前在线人数", len(WsConnAll))
	return wsConn
}

// 写入消息到队列中
func (wsConn *WsConnection) WsWrite(data []byte) error {
	select {
	case wsConn.OutChan <- data:
	case <-wsConn.CloseChan:
		return errors.New("连接已经关闭")
	}
	return nil
}

// 读取消息队列中的消息
func (wsConn *WsConnection) WsRead() ([]byte, error) {
	select {
	case msg := <-wsConn.InChan:
		// 获取到消息队列中的消息
		return msg, nil
	case <-wsConn.CloseChan:

	}
	return nil, errors.New("连接已经关闭")
}

// 关闭连接
func (wsConn *WsConnection) Close() {
	log.Println("关闭连接被调用了")
	wsConn.WsSocket.Close()
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if wsConn.isClosed == false {
		wsConn.isClosed = true
		// 删除这个连接的变量
		delete(WsConnAll, wsConn.id)
		close(wsConn.CloseChan)
	}
}

//初始玩玩家信息
func (ws *WsConnection) InitPlayer(id string) error {
	ws.Player = &manager.Player{}
	userManger, err := manager.NewUserManager(id)
	if err != nil {
		return err
	}
	ws.User = userManger

	manager.AddPlayer(ws.User.Id, ws.Player)
	registry.Register(ws.User.Id, ws.OutChan)
	return nil
}

//玩家离线
func (ws *WsConnection) OffLine(id string) {
	manager.RemoveRoom(ws.User.GetRoomId())
	manager.RemovePlayer4PvpPool(ws.User.GetPiecewise(), id)
	manager.DeletePlayer(id)
	registry.UnRegister(id)
}
