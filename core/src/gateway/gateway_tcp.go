package main

import (
	"net"
)

type TcpConn struct {
	conn             *net.TCPConn
	onlineStreamAddr string
	msgChannel       chan []byte //发送给玩家消息队列
	clientAddr       string      //client地址
	//router           *router.Router
}

func newTcpConn(conn *net.TCPConn, onlineStreamAddr string) *TcpConn {
	tcpConn := &TcpConn{
		conn:             conn,
		msgChannel:       make(chan []byte, 64),
		clientAddr:       conn.RemoteAddr().String(),
		onlineStreamAddr: onlineStreamAddr,
		//router:           router.NewRouter(),
	}
	//tcpConn.registerTcpHandler()
	return tcpConn
}

func (tcpConn *TcpConn) addMsgChannel(data []byte) {
	tcpConn.msgChannel <- data
}

func (tcpConn *TcpConn) onMessage(data []byte) {
	onlineStream := tcpServer().getOnlineStream(tcpConn.onlineStreamAddr)
	if onlineStream == nil {
		return
	}
	onlineStream.addToOnlineMsg(data)
}
