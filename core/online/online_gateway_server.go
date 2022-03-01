package main

import (
	"core/component/logger"
)

var tServer *TcpServer

type TcpServer struct {
	gatewayStreamInfos map[string]*OnlineStreamInfo //grpc流
}

func newTcpServer() *TcpServer {
	tcpServer := &TcpServer{
		gatewayStreamInfos: make(map[string]*OnlineStreamInfo),
	}
	tServer = tcpServer
	return tcpServer
}

func tcpServer() *TcpServer {
	return tServer
}

func (tcpSrv *TcpServer) getGatewayStream(addr string) *OnlineStreamInfo {
	if gatewayStream, ok := tcpSrv.gatewayStreamInfos[addr]; ok {
		return gatewayStream
	}
	return nil
}

//添加gateway grpc stream
func (tcpSrv *TcpServer) addGatewayStream(info *OnlineStreamInfo) {
	if _, ok := tcpSrv.gatewayStreamInfos[info.gatewayRemoteAddr]; ok {
		return
	}
	tcpSrv.gatewayStreamInfos[info.gatewayRemoteAddr] = info
	logger.Infof("添加grpc连接gateway--online %v", info.gatewayRemoteAddr)
}

//删除gateway grpc stream
func (tcpSrv *TcpServer) delGatewayStream(info *OnlineStreamInfo) {
	if _, ok := tcpSrv.gatewayStreamInfos[info.gatewayRemoteAddr]; !ok {
		return
	}
	delete(tcpSrv.gatewayStreamInfos, info.gatewayRemoteAddr)
	logger.Infof("删除grpc连接gateway--online %v", info.gatewayRemoteAddr)
}

