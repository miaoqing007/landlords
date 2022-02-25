package main

import (
	"bufio"
	"core/component/logger"
	"log"
	"net"
	"sync"
)

var tServer *TcpServer

type TcpServer struct {
	tcpConnects      *sync.Map //map[string]*TcpConn //tcp连接
	dialStream       *OnlineGRPCStream
	buffPool         *sync.Pool
	isOpenGRPCStream bool
}

func newTcpServer() *TcpServer {
	tcpServer := &TcpServer{
		tcpConnects: &sync.Map{},
	}
	tcpServer.initSyncPool()
	tServer = tcpServer
	return tcpServer
}

func (tcpSrv *TcpServer) getDialStream() *OnlineGRPCStream {
	return tcpSrv.dialStream
}

func (tcpSrv *TcpServer) initSyncPool() {
	tcpSrv.buffPool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024*1024)
		},
	}
}

func tcpServer() *TcpServer {
	return tServer
}

func (tcpSrv *TcpServer) addConnMsg(addr string, data []byte) {
	if value, ok := tcpSrv.tcpConnects.Load(addr); ok {
		value.(*TcpConn).msgChannel <- data
	}
}

func (tcpSrv *TcpServer) addTcpConn(conn *net.TCPConn) *TcpConn {
	if value, ok := tcpSrv.tcpConnects.Load(conn.RemoteAddr().String()); ok {
		return value.(*TcpConn)
	}
	tconn := newTcpConn(conn)
	tcpSrv.tcpConnects.Store(tconn.clientAddr, tconn)
	log.Printf("添加tcp连接 %v", conn.RemoteAddr().String())
	return tconn
}

func (tcpSrv *TcpServer) delTcpConn(addr string) {
	tcpSrv.tcpConnects.Delete(addr)
	log.Printf("移除tcp连接 %v", addr)
}

func (tcpSrv *TcpServer) getTcpConn(connAddr string) *TcpConn {
	if value, ok := tcpSrv.tcpConnects.Load(connAddr); ok {
		return value.(*TcpConn)
	}
	return nil
}

func runTcpAndGRPC() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8888")
	if err != nil {
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return
	}
	defer listener.Close()

	tcpSrv := newTcpServer()

	//连接grpc
	go runGatewayGRPCDial("127.0.0.1:9999", tcpSrv)

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			logger.Infoln("断开连接")
			return
		}
		logger.Infof("message %s->%s\n", conn.RemoteAddr(), conn.LocalAddr())
		tcoon := tcpSrv.addTcpConn(conn)
		go tcpSrv.recv(tcoon)
		go tcpSrv.send(tcoon)
	}
}

func (tcpSrv *TcpServer) recv(tconn *TcpConn) {
	ipStr := tconn.conn.RemoteAddr().String()
	defer func() {
		tconn.disTcpConnected()
		tcpSrv.delTcpConn(ipStr)
		logger.Infof("Disconnected : " + ipStr)
		tconn.conn.Close()
	}()
	reader := bufio.NewReader(tconn.conn)
	for {
		buff := tcpSrv.buffPool.Get().([]byte)
		n, err := reader.Read(buff[:])
		if err != nil {
			return
		}
		tconn.onMessage(buff[:n])
		tcpSrv.buffPool.Put(buff)
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
