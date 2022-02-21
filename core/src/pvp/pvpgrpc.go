package main

import (
	command "core/command/pb"
	"google.golang.org/grpc"
	"net"
)

type server struct{}

func runOnlinePvpGRPC() {
	g, err := net.Listen("tcp", "")
	if err != nil {
		return
	}
	s := grpc.NewServer()
	ins := new(server)
	command.RegisterOnlinePvpServer(s, ins)
	s.Serve(g)
}

func (s *server) OnlinePvpStream(streamServer command.OnlinePvp_OnlinePvpStreamServer) error {
	go send()
	go recv()
	return nil
}

func send() {

}

func recv() {

}
