package main

import (
	command "core/command/pb"
	"google.golang.org/grpc"
	"net"
)

type server struct {
}

func runGatewayOnlineGRPC() {
	g, err := net.Listen("tcp", "")
	if err != nil {
		return
	}
	s := grpc.NewServer()
	ins := new(server)
	command.RegisterGatewayOnlineServer(s, ins)
	s.Serve(g)
}

func (s *server) GatewayOnlineStream(streamServer command.GatewayOnline_GatewayOnlineStreamServer) error {
	return nil
}
