package network

import (
	"fmt"
	"net"
)

type TcpConn struct {
	Conn      net.Conn
	msgChan   chan interface{} //通知玩家消息
	msgParser *MsgParser
}

func NewTcpConn(conn net.Conn) *TcpConn {
	tcpConn := &TcpConn{
		msgParser: newMsgParser(),
		msgChan:   make(chan interface{}, 64),
		Conn:      conn,
	}
	return tcpConn
}

func (tcper *TcpConn) Read(b []byte) (int, error) { return 0, nil }

func (tcper *TcpConn) Write(b []byte) (int, error) { return 0, nil }

func (tcper *TcpConn) Run() {
	 tcper.handleRead()
	//go tcper.handleWrite()
}

func (tcper *TcpConn) handleRead() {
	for {
		msgID, msgData, err := tcper.msgParser.Read(tcper)
		if err != nil {
			return
		}
		fmt.Println(msgID, string(msgData))
	}
}

func (tcper *TcpConn) handleWrite() {
	for {
		select {
		case msg := <-tcper.msgChan:
			data, ok := msg.([]byte)
			if !ok {
				return
			}
			tcper.msgParser.Write(tcper, data...)
		}
	}
}
