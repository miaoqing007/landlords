package main

import (
	command "core/command/pb"
	"core/component/logger"
	"core/config"
	"google.golang.org/grpc"
	"io"
	"net"
	"sync"
)

var gSrv *server

type server struct {
	onlineStreams *sync.Map //map[string]*OnlineStreamInfo
}

func newServer() *server {
	srv := &server{
		onlineStreams: &sync.Map{},
	}
	gSrv = srv
	return srv
}

func gSrvGetMe() *server {
	return gSrv
}

func (srv *server) addOnlineStream(streamServer command.OnlinePvp_OnlinePvpStreamServer) *OnlineStreamInfo {
	onlineStream := newOnlineStreamInfo(streamServer)
	if _, ok := srv.onlineStreams.Load(onlineStream.onlineGRPCRemoteAddr); ok {
		return onlineStream
	}
	srv.onlineStreams.Store(onlineStream.onlineGRPCRemoteAddr, onlineStream)
	logger.Infof("添加grpc连接online--pvp %v", onlineStream.onlineGRPCRemoteAddr)
	return onlineStream
}

func (srv *server) delOnlineStream(addr string) {
	srv.onlineStreams.Delete(addr)
	logger.Infof("删除grpc连接online--pvp %v", addr)
}

func (srv *server) getOnlineStream(addr string) *OnlineStreamInfo {
	value, ok := srv.onlineStreams.Load(addr)
	if ok {
		return value.(*OnlineStreamInfo)
	}
	return nil
}

func runOnlinePvpGRPC() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+config.PvpGRPCPort)
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
	command.RegisterOnlinePvpServer(s, ins)

	logger.Infof("GRPC listening on %s", listener.Addr().String())
	s.Serve(listener)
}

func (srv *server) OnlinePvpStream(streamServer command.OnlinePvp_OnlinePvpStreamServer) error {

	streamInfo := srv.addOnlineStream(streamServer)

	go srv.send(streamInfo)
	srv.recv(streamInfo)
	return nil
}

func (srv *server) send(streamInfo *OnlineStreamInfo) {
	for {
		select {
		case msg := <-streamInfo.msgToOnline:
			err := streamInfo.onlineStream.Send(msg)
			if err != nil {
				return
			}
		}
	}
}

func (srv *server) recv(streamInfo *OnlineStreamInfo) {
	defer func() {
		srv.delOnlineStream(streamInfo.onlineGRPCRemoteAddr)
	}()
	for {
		out, err := streamInfo.onlineStream.Recv()
		if err != nil || err == io.EOF {
			return
		}
		out.Addr = streamInfo.onlineGRPCRemoteAddr
		WorldGetMe().addMsgChannel(out)
	}
}
