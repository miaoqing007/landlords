package websocket

import (
	"github.com/gorilla/websocket"
	"landlords/client_handler"
	"landlords/helper/stack"
	"landlords/misc/packet"
	. "landlords/obj"
	. "landlords/wsconnection"
	"log"
	"net/http"
	"time"
)

const (
	// 允许等待的写入时间
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// 最大的连接ID，每次连接都加1 处理
var maxConnId int64

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有的CORS 跨域请求，正式环境可以关闭
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(resp http.ResponseWriter, req *http.Request) {
	// 应答客户端告知升级连接为websocket
	wsSocket, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Println("升级为websocket失败", err.Error())
		return
	}
	maxConnId++
	// TODO 如果要控制连接数可以计算，wsConnAll长度
	// 连接数保持一定数量，超过的部分不提供服务
	wsConn := NewWsConnection(wsSocket, maxConnId)

	// 处理器,发送定时信息，避免意外关闭
	go processLoop(wsConn)
	// 读协程
	go wsReadLoop(wsConn)
	// 写协程
	go wsWriteLoop(wsConn)
}

// 处理队列中的消息
func processLoop(wsConn *WsConnection) {
	// 处理消息队列中的消息
	// 获取到消息队列中的消息，处理完成后，发送消息给客户端
	for {
		msg, err := wsConn.WsRead()
		if err != nil {
			log.Println("获取消息出现错误", err.Error())
			break
		}
		log.Println("接收到消息", string(msg))

		reader := packet.Reader(msg)
		c, _ := reader.ReadS16()
		byts := executeHandler(int16(c), wsConn, reader)
		for _, byt := range byts {
			err = wsConn.WsWrite(byt)
			if err != nil {
				log.Println("发送消息给客户端出现错误", err.Error())
				break
			}
		}
	}
}

// 处理消息队列中的消息
func wsReadLoop(wsConn *WsConnection) {
	// 设置消息的最大长度
	wsConn.WsSocket.SetReadLimit(maxMessageSize)
	wsConn.WsSocket.SetReadDeadline(time.Now().Add(pongWait))
	for {
		// 读一个message
		_, data, err := wsConn.WsSocket.ReadMessage()
		if err != nil {
			websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure)
			log.Println("消息读取出现错误", err.Error())
			wsConn.Close()
			return
		}
		//req := &WsMessage{
		//	msgType,
		//	data,
		//}
		// 放入请求队列,消息入栈
		select {
		case wsConn.InChan <- data:
		case <-wsConn.CloseChan:
			return
		}
	}
}

// 发送消息给客户端
func wsWriteLoop(wsConn *WsConnection) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		// 取一个应答
		case msg := <-wsConn.OutChan:
			// 写给websocket
			if err := wsConn.WsSocket.WriteMessage(0, msg); err != nil {
				log.Println("发送消息给客户端发生错误", err.Error())
				// 切断服务
				wsConn.Close()
				return
			}
		case <-wsConn.CloseChan:
			// 获取到关闭通知
			return
		case <-ticker.C:
			// 出现超时情况
			wsConn.WsSocket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := wsConn.WsSocket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 启动程序
func StartWebSocket(addrPort string) {
	WsConnAll = make(map[int64]*WsConnection)
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(addrPort, nil)
}

func executeHandler(code int16, ws *WsConnection, reader *packet.Packet) [][]byte {
	defer stack.PrintRecoverFromPanic()
	handle := client_handler.Handlers[code]
	if handle == nil {
		return nil
	}
	retByte := handle(ws, reader)
	return retByte
}
