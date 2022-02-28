package main

import (
	"bufio"
	"core/component/logger"
	"core/config"
	"io"
	"log"
	"net"
	"sync"
)

var tServer *TcpServer

type TcpServer struct {
	tcpConnects      *sync.Map  //map[string]*TcpConn //tcp连接
	buffPool         *sync.Pool //tcpconn buff池 make([]byte,1024*1024)
	dialStream       *OnlineGRPCStream
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
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+config.GatewayTCPPort)
	if err != nil {
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return
	}
	defer listener.Close()
	logger.Infof("TCP listening on %s", listener.Addr().String())

	tcpSrv := newTcpServer()

	//连接grpc
	go runGatewayGRPCDial(tcpSrv)

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			//logger.Infoln("断开连接")
			return
		}
		logger.Infof("message %s->%s\n", conn.RemoteAddr(), conn.LocalAddr())
		tcoon := tcpSrv.addTcpConn(conn)
		go tcpSrv.send(tcoon)
		go tcpSrv.recv(tcoon)
	}
}

func (tcpSrv *TcpServer) recv(tconn *TcpConn) {
	ipStr := tconn.conn.RemoteAddr().String()
	defer func() {
		tconn.disTcpConnected()
		tcpSrv.delTcpConn(ipStr)
		tconn.closeWriteChannel <- true
		tconn.conn.Close()
		logger.Infof("Disconnected : " + ipStr)
	}()
	reader := bufio.NewReader(tconn.conn)
	for {
		buff := tcpSrv.buffPool.Get().([]byte)
		n, err := reader.Read(buff[:])
		if err != nil || err == io.EOF {
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
			if _, err := tconn.conn.Write(data); err != nil {
				return
			}
		case <-tconn.closeWriteChannel:
			logger.Infof("关闭玩家(%v) conn写入连接", tconn.playerId)
			return
		}
	}
}
