package wsconnection

import (
	"errors"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"landlords/manager"
	. "landlords/obj"
	"landlords/registry"
	"sync"
)

// ws 的所有连接
// 用于广播
var (
	//WsConnAll = make(map[int]*WsConnection)
	currentConn int
)

// 客户端连接
type WsConnection struct {
	WsSocket *websocket.Conn // 底层websocket
	InChan   chan *WsMessage // 读队列
	OutChan  chan *WsMessage // 写队列

	mutex     sync.Mutex // 避免重复关闭管道,加锁处理
	isClosed  bool
	CloseChan chan byte // 关闭通知
	*manager.Player
}

func NewWsConnection(wsSocket *websocket.Conn) *WsConnection {
	wsConn := &WsConnection{
		WsSocket:  wsSocket,
		InChan:    make(chan *WsMessage, 1000),
		OutChan:   make(chan *WsMessage, 1000),
		CloseChan: make(chan byte),
		isClosed:  false,
	}
	currentConn++
	glog.Info("当前连接数", currentConn)
	return wsConn
}

// 写入消息到队列中
func (wsConn *WsConnection) WsWrite(data *WsMessage) error {
	select {
	case wsConn.OutChan <- data:
	case <-wsConn.CloseChan:
		return errors.New("连接已经关闭")
	}
	return nil
}

// 读取消息队列中的消息
func (wsConn *WsConnection) WsRead() (*WsMessage, error) {
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
	glog.Info("关闭连接被调用了")
	wsConn.WsSocket.Close()
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if wsConn.isClosed == false {
		wsConn.isClosed = true
		// 删除这个连接的变量
		currentConn--
		close(wsConn.CloseChan)
		if wsConn.Player != nil && wsConn.User != nil {
			wsConn.OffLine(wsConn.User.Id)
		}
	}
}

//初始玩玩家信息
func (ws *WsConnection) InitPlayer(account, password string) error {
	ws.Player = &manager.Player{}
	userManger, err := manager.NewUserManager(account, password)
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
