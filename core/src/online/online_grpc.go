package main

import (
	command "core/command/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"io"
	"net"
)

type server struct {
	tcpSrv *TcpServer
}

func newServer() *server {
	return &server{
		tcpSrv: newTcpServer(),
	}
}

type OnlineStreamInfo struct {
	onlineStream command.GatewayOnline_GatewayOnlineStreamServer
	toGatewayMsgChan chan *command.ServerPlayerMsgData //online-->gateway-->client
	addr             string                            //
}

func newGatewayInfo(streamServer command.GatewayOnline_GatewayOnlineStreamServer) *OnlineStreamInfo {
	gateway := &OnlineStreamInfo{
		toGatewayMsgChan: make(chan *command.ServerPlayerMsgData, 1024),
		onlineStream:     streamServer,
	}
	return gateway
}

func (g *OnlineStreamInfo) addToGatewayMsg(data []byte, clientAddr string) {
	g.toGatewayMsgChan <- &command.ServerPlayerMsgData{Data: data, ClientAddr: clientAddr}
}


func (g *OnlineStreamInfo) getRemoteAddr() string {
	pr, ok := peer.FromContext(g.onlineStream.Context())
	if !ok {
		return ""
	}
	if pr.Addr == net.Addr(nil) {
		return ""
	}
	return pr.Addr.String()
}

func runGatewayOnlineGRPC() {
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
	ins := newServer()
	command.RegisterGatewayOnlineServer(s, ins)
	s.Serve(listener)
}

func (s *server) GatewayOnlineStream(streamServer command.GatewayOnline_GatewayOnlineStreamServer) error {
	gatewayInfo := newGatewayInfo(streamServer)
	s.tcpSrv.addOnlineStream(gatewayInfo)
	go s.send(gatewayInfo)
	s.recv(gatewayInfo)
	return nil
}

func (s *server) send(onlineStream *OnlineStreamInfo) {
	for {
		out, ok := <-onlineStream.toGatewayMsgChan
		if !ok {
			return
		}
		onlineStream.onlineStream.Send(out)
	}
}

func (s *server) recv(onlineStreamInfo *OnlineStreamInfo) {
	defer func() {
		s.tcpSrv.delOnlineStream(onlineStreamInfo)
	}()
	for {
		out, err := onlineStreamInfo.onlineStream.Recv()
		if err == io.EOF || err != nil {
			return
		}
		WorldGetMe().addPlayer(out.PlayerId, out.ClientAddr, onlineStreamInfo.getRemoteAddr())
		WorldGetMe().sendFromGatewayMsgChan(out)
	}
}
