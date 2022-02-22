package main

import (
	command "core/command/pb"
	"core/component/router"
	"google.golang.org/grpc/peer"
	"log"
	"net"
	"sync"
)

type GatewayInfo struct {
	onlineStream      command.GatewayOnline_GatewayOnlineStreamServer
	recvClientMsgChan chan *command.ClientPlayerMsgData
	sendClientMsgChan chan *command.ClientPlayerMsgData
}

func newGatewayInfo(streamServer command.GatewayOnline_GatewayOnlineStreamServer) *GatewayInfo {
	gateway := &GatewayInfo{
		recvClientMsgChan: make(chan *command.ClientPlayerMsgData, 64),
		sendClientMsgChan: make(chan *command.ClientPlayerMsgData, 64),
		onlineStream:      streamServer,
	}
	return gateway
}

func (g *GatewayInfo) addSendClientMsgChan(data []byte) {
	g.sendClientMsgChan <- &command.ClientPlayerMsgData{Data: data}
}

func (g *GatewayInfo) getRemoteAddr() string {
	pr, ok := peer.FromContext(g.onlineStream.Context())
	if !ok {
		return ""
	}
	if pr.Addr == net.Addr(nil) {
		return ""
	}
	return pr.Addr.String()
}

type TcpServer struct {
	tcpConnects   map[string]*TcpConn     //tcp连接
	onlineStreams map[string]*GatewayInfo //grpc流
	router        *router.Router
	sync.RWMutex
}

type TcpConn struct {
	conn       *net.TCPConn
	msgChannel chan []byte //发送给玩家消息队列
}

func newTcpConn(conn *net.TCPConn) *TcpConn {
	tcpConn := &TcpConn{
		conn:       conn,
		msgChannel: make(chan []byte, 64),
	}
	return tcpConn
}

func newTcpServer() *TcpServer {
	tcpServer := &TcpServer{
		tcpConnects:   make(map[string]*TcpConn),
		onlineStreams: make(map[string]*GatewayInfo),
		router:        router.NewRouter(),
	}
	tcpServer.registerGatewayOnline()
	return tcpServer
}

func (tcpSrv *TcpServer) addConnMsg(addr string, data []byte) {
	if conn, ok := tcpSrv.tcpConnects[addr]; ok {
		conn.msgChannel <- data
	}
}

//添加online grpc stream
func (tcpSrv *TcpServer) addOnlineStream(info *GatewayInfo) {
	if _, ok := tcpSrv.onlineStreams[info.getRemoteAddr()]; ok {
		return
	}
	tcpSrv.onlineStreams[info.getRemoteAddr()] = info
}

//删除online grpc stream
func (tcpSrv *TcpServer) delOnlineStream(info *GatewayInfo) {
	if _, ok := tcpSrv.onlineStreams[info.getRemoteAddr()]; !ok {
		return
	}
	delete(tcpSrv.onlineStreams, info.getRemoteAddr())
}

func (tcpSrv *TcpServer) addTcpConn(conn *net.TCPConn) {
	if _, ok := tcpSrv.tcpConnects[conn.RemoteAddr().String()]; ok {
		return
	}
	tcpSrv.RLock()
	tcpSrv.tcpConnects[conn.RemoteAddr().String()] = newTcpConn(conn)
	tcpSrv.RUnlock()
	log.Printf("添加tcp连接 %v", conn.RemoteAddr().String())
}

func (tcpSrv *TcpServer) delTcpConn(addr string) {
	tcpSrv.RLock()
	delete(tcpSrv.tcpConnects, addr)
	tcpSrv.RUnlock()
	log.Printf("移除tcp连接 %v", addr)
}
