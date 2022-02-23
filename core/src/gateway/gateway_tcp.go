package main

import (
	"net"
)

type TcpConn struct {
	conn       *net.TCPConn
	msgChannel chan []byte //发送给玩家消息队列
	clientAddr string      //client地址
}

func newTcpConn(conn *net.TCPConn) *TcpConn {
	tcpConn := &TcpConn{
		conn:       conn,
		msgChannel: make(chan []byte, 64),
		clientAddr: conn.RemoteAddr().String(),
	}
	return tcpConn
}

func (tcpConn *TcpConn) addMsgChannel(data []byte) {
	tcpConn.msgChannel <- data
}

func (tcpConn *TcpConn) onMessage(data []byte) {
	onlineStream := tcpServer().dialStream
	if onlineStream == nil {
		return
	}
	onlineStream.addMsgChannel(data, tcpConn.clientAddr)
}
