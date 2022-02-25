package main

import (
	command "core/command/pb"
	"core/component/logger"
	"core/component/router"
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
	gatewayStream     command.GatewayOnline_GatewayOnlineStreamServer
	toGatewayMsgChan  chan *command.ServerPlayerMsgData //online-->gateway-->client
	gatewayRemoteAddr string                            //
	router            *router.Router
}

func newGatewayInfo(streamServer command.GatewayOnline_GatewayOnlineStreamServer) *OnlineStreamInfo {
	gateway := &OnlineStreamInfo{
		toGatewayMsgChan: make(chan *command.ServerPlayerMsgData, 1024),
		gatewayStream:    streamServer,
		router:           router.NewRouter(),
	}
	gateway.gatewayRemoteAddr = gateway.getRemoteAddr()
	gateway.registerGatewayOnlineHandler()
	return gateway
}

//gatewayOnline消息注册
func (os *OnlineStreamInfo) registerGatewayOnlineHandler() {
	os.router.RegisterGatewayOnline(uint16(command.Command_ClientInOnline), os.clientInOnlineHandler)
	os.router.RegisterGatewayOnline(uint16(command.Command_ClientOutOnline), os.clientOutOnlineHandler)
}

func (os *OnlineStreamInfo) addToGatewayMsg(data []byte, clientAddr string) {
	os.toGatewayMsgChan <- &command.ServerPlayerMsgData{Data: data, ClientAddr: clientAddr}
}

func (os *OnlineStreamInfo) getRemoteAddr() string {
	pr, ok := peer.FromContext(os.gatewayStream.Context())
	if !ok {
		return ""
	}
	if pr.Addr == net.Addr(nil) {
		return ""
	}
	return pr.Addr.String()
}

func (os *OnlineStreamInfo) onMessage(out *command.ClientPlayerMsgData) {
	if _, err := os.router.RouterGatewayOnlineMsg(out.Data); err != nil {
		//转发到其他服玩家数据
		WorldGetMe().sendFromOtherServerMsgChan(out.PlayerId, out.Data)
		logger.Infof("转发到其他服玩家(%v)数据", out.PlayerId)
	}
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

	go runOnlinePvpGRPC("127.0.0.1:2222")

	s := grpc.NewServer()
	ins := newServer()
	command.RegisterGatewayOnlineServer(s, ins)
	s.Serve(listener)
}

func (s *server) GatewayOnlineStream(streamServer command.GatewayOnline_GatewayOnlineStreamServer) error {
	gatewayInfo := newGatewayInfo(streamServer)
	s.tcpSrv.addGatewayStream(gatewayInfo)
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
		if err := onlineStream.gatewayStream.Send(out); err != nil {
			return
		}
	}
}

func (s *server) recv(onlineStreamInfo *OnlineStreamInfo) {
	defer func() {
		s.tcpSrv.delGatewayStream(onlineStreamInfo)
	}()
	for {
		out, err := onlineStreamInfo.gatewayStream.Recv()
		if err == io.EOF || err != nil {
			return
		}
		onlineStreamInfo.onMessage(out)
		//WorldGetMe().addPlayer(out.PlayerId, out.ClientAddr, onlineStreamInfo.gatewayRemoteAddr)
		//WorldGetMe().sendFromGatewayMsgChan(out.PlayerId, out.Data, out.ClientAddr)
	}
}
