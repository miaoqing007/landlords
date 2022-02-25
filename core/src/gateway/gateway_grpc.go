package main

import (
	"context"
	command "core/command/pb"
	"core/component/logger"
	"google.golang.org/grpc"
	"io"
	"time"
)

type OnlineGRPCStream struct {
	client             command.GatewayOnlineClient
	msg2Online         chan *command.ClientPlayerMsgData //gateway-->online
	isSuccessConnected bool                              //是否已成功连接
}

func newOnlineGRPCStream(client command.GatewayOnlineClient) *OnlineGRPCStream {
	gs := &OnlineGRPCStream{
		client:     client,
		msg2Online: make(chan *command.ClientPlayerMsgData, 1024),
	}
	return gs
}

func (gs *OnlineGRPCStream) addMsgChannel(data []byte, clientAddr string) {
	gs.msg2Online <- &command.ClientPlayerMsgData{PlayerId: 123456, Data: data, ClientAddr: clientAddr}
}

func (gs *OnlineGRPCStream) openStream() {
	stream, err := gs.client.GatewayOnlineStream(context.Background())
	if err != nil {
		return
	}
	gs.isSuccessConnected = true
	logger.Info("开启grpc流成功 gateway-->online")
	defer func() {
		gs.isSuccessConnected = false
		stream.CloseSend()
		logger.Info("grpc流断开连接 gateway-->online")
	}()
	go func() {
		for {
			select {
			case msg := <-gs.msg2Online:
				if err := stream.Send(msg); err != nil {
					return
				}
			}
		}
	}()
	for {
		msg, err := stream.Recv()
		if err == io.EOF || err != nil {
			return
		}
		gs.onMessage(msg)
	}
}

func (gs *OnlineGRPCStream) onMessage(msg *command.ServerPlayerMsgData) {
	conn := tcpServer().getTcpConn(msg.ClientAddr)
	if conn != nil {
		conn.addMsgChannel(msg.Data)
	}
}

func runGatewayGRPCDial(addr string, tcpSrv *TcpServer) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logger.Error(err)
		return
	}
	client := command.NewGatewayOnlineClient(conn)

	gs := newOnlineGRPCStream(client)
	tcpSrv.dialStream = gs
	for {
		if !gs.isSuccessConnected {
			gs.openStream()
		}
		time.Sleep(5 * time.Second)
	}
}
