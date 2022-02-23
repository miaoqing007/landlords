package main

import (
	"context"
	command "core/command/pb"
	"core/component/router"
	"fmt"
	"google.golang.org/grpc"
	"io"
)

type GRPCStream struct {
	client       command.GatewayOnlineClient
	msgChannel   chan *command.ClientPlayerMsgData
	recvChannel  chan *command.ServerPlayerMsgData
	closeChannel chan bool
	router       *router.Router
}

func newGRPCStream(client command.GatewayOnlineClient) *GRPCStream {
	gs := &GRPCStream{
		client:       client,
		msgChannel:   make(chan *command.ClientPlayerMsgData, 1024),
		recvChannel:  make(chan *command.ServerPlayerMsgData, 1024),
		closeChannel: make(chan bool, 0),
		router:       router.NewRouter(),
	}
	return gs
}

func (gs *GRPCStream) close() {
	gs.closeChannel <- true
}

func (gs *GRPCStream) addMsgChannel(data []byte, clientAddr string) {
	gs.msgChannel <- &command.ClientPlayerMsgData{PlayerId: 123456, Data: data, ClientAddr: clientAddr}
}

func (gs *GRPCStream) openStream() {
	stream, err := gs.client.GatewayOnlineStream(context.Background())
	if err != nil {
		return
	}
	defer func() {
		gs.close()
	}()
	go func() {
		for {
			select {
			case msg := <-gs.msgChannel:
				stream.Send(msg)
			}
		}
	}()
	for {
		msg, err := stream.Recv()
		if err == io.EOF || err != nil {
			return
		}
		conn := tcpServer().getTcpConn(msg.ClientAddr)
		if conn != nil {
			conn.addMsgChannel(msg.Data)
		}
	}
}

func runGRPCDial(addr string, tcpSrv *TcpServer) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	client := command.NewGatewayOnlineClient(conn)

	gs := newGRPCStream(client)
	tcpSrv.dialStream = gs

	gs.openStream()
}
