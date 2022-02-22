package main

import (
	command "core/command/pb"
	"google.golang.org/grpc"
	"io"
	"net"
)

type server struct {
	tcpServer *TcpServer
}

func newServer(tcpServer *TcpServer) *server {
	srv := &server{
		tcpServer: tcpServer,
	}
	return srv
}

func runGatewayOnlineGRPC(tcpSrv *TcpServer) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9999")
	if err != nil {
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return
	}
	defer listener.Close()
	s := grpc.NewServer()
	ins := newServer(tcpSrv)
	command.RegisterGatewayOnlineServer(s, ins)
	s.Serve(listener)
}

func (s *server) GatewayOnlineStream(streamServer command.GatewayOnline_GatewayOnlineStreamServer) error {

	gatewayInfo := newGatewayInfo(streamServer)
	s.tcpServer.addOnlineStream(gatewayInfo)

	go s.send(gatewayInfo)
	s.recv(gatewayInfo)
	return nil
}

func (s *server) send(gatewayInfo *GatewayInfo) {
	for {
		out, ok := <-gatewayInfo.sendClientMsgChan
		if !ok {
			return
		}
		gatewayInfo.onlineStream.Send(out)
	}
}

func (s *server) recv(gatewayInfo *GatewayInfo) {
	defer func() {
		s.tcpServer.delOnlineStream(gatewayInfo)
	}()
	for {
		out, err := gatewayInfo.onlineStream.Recv()
		if err == io.EOF || err != nil {
			return
		}
		s.tcpServer.addConnMsg(gatewayInfo.getRemoteAddr(), out.Data)
	}
}
