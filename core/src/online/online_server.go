package main

import (
	"log"
	"sync"
)

var tServer *TcpServer

type TcpServer struct {
	onlineGatewayInfos map[string]*OnlineStreamInfo //grpc流
	sync.RWMutex
}

func newTcpServer() *TcpServer {
	tcpServer := &TcpServer{
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

//添加online grpc stream
func (tcpSrv *TcpServer) addOnlineStream(info *OnlineStreamInfo) {
	if _, ok := tcpSrv.onlineGatewayInfos[info.getRemoteAddr()]; ok {
		return
	}
	tcpSrv.onlineGatewayInfos[info.getRemoteAddr()] = info
	log.Printf("添加grpc连接gateway--online %v", info.getRemoteAddr())
}

//删除online grpc stream
func (tcpSrv *TcpServer) delOnlineStream(info *OnlineStreamInfo) {
	if _, ok := tcpSrv.onlineGatewayInfos[info.getRemoteAddr()]; !ok {
		return
	}
	delete(tcpSrv.onlineGatewayInfos, info.getRemoteAddr())
	log.Printf("删除grpc连接gateway--online %v", info.getRemoteAddr())
}

func (tcpSrv *TcpServer) randAnOnlineStreamAddr() string {
	for _, onlineStream := range tcpSrv.onlineGatewayInfos {
		return onlineStream.getRemoteAddr()
	}
	return ""
}

