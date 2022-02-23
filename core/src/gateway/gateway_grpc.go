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

func newServer(tcpSrv *TcpServer) *server {
	return &server{
		tcpSrv: tcpSrv,
	}
}

type OnlineStreamInfo struct {
	onlineStream command.GatewayOnline_GatewayOnlineStreamServer
	//toClientMsgChan chan []byte                       //online-->grpc-->client
	toOnlineMsgChan chan *command.ClientPlayerMsgData //client-->grpc-->online
	addr            string                            //
	//router          *router.Router
}

func newGatewayInfo(streamServer command.GatewayOnline_GatewayOnlineStreamServer) *OnlineStreamInfo {
	gateway := &OnlineStreamInfo{
		//toClientMsgChan: make(chan []byte, 1024),
		toOnlineMsgChan: make(chan *command.ClientPlayerMsgData, 1024),
		onlineStream:    streamServer,
		//router:          router.NewRouter(),
	}
	return gateway
}

func (g *OnlineStreamInfo) addToOnlineMsg(data []byte) {
	g.toOnlineMsgChan <- &command.ClientPlayerMsgData{Data: data}
}

//func (g *OnlineStreamInfo) addToClientMsg(data []byte) {
//	g.toClientMsgChan <- data
//}

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
	s.tcpSrv.addOnlineStream(gatewayInfo)
	go s.send(gatewayInfo)
	s.recv(gatewayInfo)
	return nil
}

func (s *server) send(onlineStream *OnlineStreamInfo) {
	for {
		out, ok := <-onlineStream.toOnlineMsgChan
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
		tconn := s.tcpSrv.getTcpConn(out.ClientAddr)
		if tconn != nil {
			tconn.addMsgChannel(out.Data)
		}
	}
}
