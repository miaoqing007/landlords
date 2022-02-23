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
	msgChannel   chan *command.ServerPlayerMsgData
	recvChannel  chan *command.ClientPlayerMsgData
	closeChannel chan bool
	router       *router.Router
}

func newGRPCStream(client command.GatewayOnlineClient) *GRPCStream {
	gs := &GRPCStream{
		client:       client,
		msgChannel:   make(chan *command.ServerPlayerMsgData, 1024),
		recvChannel:  make(chan *command.ClientPlayerMsgData, 1024),
		closeChannel: make(chan bool, 0),
		router:       router.NewRouter(),
	}
	go gs.loop()
	return gs
}

func (gs *GRPCStream) close() {
	gs.closeChannel <- true
}

func (gs *GRPCStream) loop() {
	for {
		select {
		case msg := <-gs.recvChannel:
			WorldGetMe().sendFromGatewayMsgChan(msg)
		case <-gs.closeChannel:
			return
		}
	}
}

func (gs *GRPCStream) addRecvChannel(msg *command.ClientPlayerMsgData) {
	gs.recvChannel <- msg
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
		gs.addRecvChannel(msg)
	}
}

func runGRPCDial(addr string) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	client := command.NewGatewayOnlineClient(conn)

	gs := newGRPCStream(client)

	gs.openStream()
}
