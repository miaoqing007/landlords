package main

import (
	"bufio"
	command "core/command/pb"
	"fmt"
	"github.com/golang/glog"
	"log"
	"net"
	"sync"
)

var tServer *TcpServer

type TcpServer struct {
	tcpConnects map[string]*TcpConn //tcp连接
	dialStream  *GRPCStream
	sync.RWMutex
}

func newTcpServer() *TcpServer {
	tcpServer := &TcpServer{
		tcpConnects: make(map[string]*TcpConn),
	}
	tServer = tcpServer
	return tcpServer
}

func tcpServer() *TcpServer {
	return tServer
}

func (tcpSrv *TcpServer) addConnMsg(addr string, data []byte) {
	if conn, ok := tcpSrv.tcpConnects[addr]; ok {
		conn.msgChannel <- data
	}
}

func (tcpSrv *TcpServer) addTcpConn(conn *net.TCPConn) *TcpConn {
	if tconn, ok := tcpSrv.tcpConnects[conn.RemoteAddr().String()]; ok {
		return tconn
	}
	tconn := newTcpConn(conn)
	tcpSrv.RLock()
	tcpSrv.tcpConnects[tconn.clientAddr] = tconn
	tcpSrv.RUnlock()
	log.Printf("添加tcp连接 %v", conn.RemoteAddr().String())
	return tconn
}

func (tcpSrv *TcpServer) delTcpConn(addr string) {
	tcpSrv.RLock()
	delete(tcpSrv.tcpConnects, addr)
	tcpSrv.RUnlock()
	log.Printf("移除tcp连接 %v", addr)
}

func (tcpSrv *TcpServer) addRecvChannel(msg *command.ServerPlayerMsgData) {
	tcpSrv.dialStream.recvChannel <- msg
}

func (tcpSrv *TcpServer) getTcpConn(connAddr string) *TcpConn {
	if conn, ok := tcpSrv.tcpConnects[connAddr]; ok {
		return conn
	}
	return nil
}

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
	go runGRPCDial("127.0.0.1:9999", tcpSrv)

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			glog.Infoln("断开连接")
			return
		}
		glog.Infof("message %s->%s\n", conn.RemoteAddr(), conn.LocalAddr())
		tcoon := tcpSrv.addTcpConn(conn)
		go tcpSrv.recv(tcoon)
		go tcpSrv.send(tcoon)
	}
}

func (tcpSrv *TcpServer) recv(tconn *TcpConn) {
	ipStr := tconn.conn.RemoteAddr().String()
	defer func() {
		tcpSrv.delTcpConn(ipStr)
		fmt.Println("Disconnected : " + ipStr)
		tconn.conn.Close()
	}()
	reader := bufio.NewReader(tconn.conn)
	for {
		buff := make([]byte, 1024*1024)
		n, err := reader.Read(buff[:])
		if err != nil {
			return
		}
		tconn.onMessage(buff[:n])
	}
}

func (tcpSrv *TcpServer) send(tconn *TcpConn) {
	for {
		select {
		case data := <-tconn.msgChannel:
			tconn.conn.Write(data)
		}
	}
}
