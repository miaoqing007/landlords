package main

import (
	"bufio"
	"fmt"
	"github.com/golang/glog"
	"net"
)

var ch = make(chan []byte, 10)

func run(addr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return
	}
	defer listener.Close()

	tcpSrv := newTcpServer()

	//grpc
	go runGatewayOnlineGRPC(tcpSrv)

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			glog.Infoln("断开连接")
			return
		}
		glog.Infof("message %s->%s\n", conn.RemoteAddr(), conn.LocalAddr())
		tcpSrv.addTcpConn(conn)
		go tcpSrv.recv(conn)
		go tcpSrv.send(conn)
	}
}

func (tcpSrv *TcpServer) recv(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	defer func() {
		tcpSrv.delTcpConn(ipStr)
		fmt.Println("Disconnected : " + ipStr)
		conn.Close()
	}()
	reader := bufio.NewReader(conn)
	for {
		buff := make([]byte, 1024*1024)
		n, err := reader.Read(buff[:])
		if err != nil {
			return
		}
		tcpSrv.router.RouterInnerMsg(buff[:n])
	}
}

func (tcpSrv *TcpServer) send(conn *net.TCPConn) {
	for {
		select {
		case data := <-tcpSrv.tcpConnects[conn.RemoteAddr().String()].msgChannel:
			conn.Write(data)
		}
	}
}

func OnMessage() {

}
