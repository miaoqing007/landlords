package main

import (
	"bufio"
	"fmt"
	"github.com/golang/glog"
	"log"
	"net"
	"sync"
)

var tServer *TcpServer

type TcpServer struct {
	tcpConnects        map[string]*TcpConn          //tcp连接
	onlineGatewayInfos map[string]*OnlineStreamInfo //grpc流
	sync.RWMutex
}

func newTcpServer() *TcpServer {
	tcpServer := &TcpServer{
		tcpConnects:        make(map[string]*TcpConn),
		onlineGatewayInfos: make(map[string]*OnlineStreamInfo),
	}
	tServer = tcpServer
	return tcpServer
}

func tcpServer() *TcpServer {
	return tServer
}

func (tcpSrv *TcpServer) getOnlineStream(addr string) *OnlineStreamInfo {
	if gatewayInfo, ok := tcpSrv.onlineGatewayInfos[addr]; ok {
		return gatewayInfo
	}
	return nil
}

func (tcpSrv *TcpServer) addConnMsg(addr string, data []byte) {
	if conn, ok := tcpSrv.tcpConnects[addr]; ok {
		conn.msgChannel <- data
	}
}

//添加online grpc stream
func (tcpSrv *TcpServer) addOnlineStream(info *OnlineStreamInfo) {
	if _, ok := tcpSrv.onlineGatewayInfos[info.getRemoteAddr()]; ok {
		return
	}
	tcpSrv.onlineGatewayInfos[info.getRemoteAddr()] = info
	log.Printf("添加grpc连接online--gateway %v", info.getRemoteAddr())
}

//删除online grpc stream
func (tcpSrv *TcpServer) delOnlineStream(info *OnlineStreamInfo) {
	if _, ok := tcpSrv.onlineGatewayInfos[info.getRemoteAddr()]; !ok {
		return
	}
	delete(tcpSrv.onlineGatewayInfos, info.getRemoteAddr())
	log.Printf("删除grpc连接online--gateway %v", info.getRemoteAddr())
}

func (tcpSrv *TcpServer) addTcpConn(conn *net.TCPConn) *TcpConn {
	if tconn, ok := tcpSrv.tcpConnects[conn.RemoteAddr().String()]; ok {
		return tconn
	}
	tconn := newTcpConn(conn, tcpSrv.randAnOnlineStreamAddr())
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

func (tcpSrv *TcpServer) randAnOnlineStreamAddr() string {
	for _, onlineStream := range tcpSrv.onlineGatewayInfos {
		return onlineStream.getRemoteAddr()
	}
	return ""
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
	go runGatewayOnlineGRPC(tcpSrv)

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
