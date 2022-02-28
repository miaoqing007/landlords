package main

import (
	command "core/command/pb"
	"core/component/logger"
	"core/component/router"
	"net"
)

type TcpConn struct {
	conn              *net.TCPConn
	msgChannel        chan []byte //发送给玩家消息队列
	clientAddr        string      //client地址
	playerId          uint64      //玩家id
	closeWriteChannel chan bool   //
	router            *router.Router
}

func newTcpConn(conn *net.TCPConn) *TcpConn {
	tcpConn := &TcpConn{
		conn:              conn,
		msgChannel:        make(chan []byte, 64),
		clientAddr:        conn.RemoteAddr().String(),
		router:            router.NewRouter(),
		closeWriteChannel: make(chan bool, 0),
	}
	tcpConn.registerGatewayOnlineHandler()
	return tcpConn
}

//gatewayOnline消息注册
func (tcpConn *TcpConn) registerGatewayOnlineHandler() {
	tcpConn.router.RegisterGatewayOnline(uint16(command.Command_ClientInOnline), tcpConn.clientInOnlineHandler)
	//client主动退出游戏
	tcpConn.router.RegisterGatewayOnline(uint16(command.Command_ClientOutOnline), tcpConn.clientOutOnlineHandler)
}

func (tcpConn *TcpConn) addMsgChannel(data []byte) {
	tcpConn.msgChannel <- data
}

func (TcpConn *TcpConn) bindPlayerId(playerId uint64) {
	TcpConn.playerId = playerId
}

func (tcpConn *TcpConn) disTcpConnected() {
	if tcpConn.playerId != 0 {
		msgSend := &command.ClientOutOnline_Online{
			PlayerId: tcpConn.playerId,
		}
		tcpConn.sendMsgToOnline(command.Command_ClientOutOnline, msgSend)
		logger.Infof("玩家(%v)断开online", tcpConn.playerId)
	}
}

func (tcpConn *TcpConn) sendMsgToOnline(cmd command.Command, msg interface{}) {
	data, err := tcpConn.router.Marshal(uint16(cmd), msg)
	if err != nil {
		return
	}
	stream := tcpServer().getDialStream()
	if stream != nil {
		stream.addMsgChannel(data, tcpConn.clientAddr)
	}
}

func (tcpConn *TcpConn) onMessage(data []byte) {
	onlineStream := tcpServer().getDialStream()
	if onlineStream == nil {
		return
	}
	if _, err := tcpConn.router.RouterGatewayOnlineMsg(data); err == nil {
		return
	}
	onlineStream.addMsgChannel(data, tcpConn.clientAddr)
}
