package main

import (
	"context"
	command "core/command/pb"
	"core/component/logger"
	"google.golang.org/grpc"
	"io"
	"time"
)

var pStream *PvpGRPCStream

type PvpGRPCStream struct {
	client             command.OnlinePvpClient      //pvpGrpcClient
	msgPvpChannel      chan *command.Online2PvpInfo //online-->pvp
	isSuccessConnected bool                         //是否已成功连接
}

func newPvpGRPCStream(client command.OnlinePvpClient) *PvpGRPCStream {
	ps := &PvpGRPCStream{
		client:        client,
		msgPvpChannel: make(chan *command.Online2PvpInfo, 1024),
	}
	pStream = ps
	return ps
}

func PvpStreamGetMe() *PvpGRPCStream {
	return pStream
}

func (ps *PvpGRPCStream) addPvpMsgChannel(playerId uint64, data []byte) {
	ps.msgPvpChannel <- &command.Online2PvpInfo{PlayerId: playerId, Data: data}
}

func runOnlinePvpGRPC(addr string) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return
	}
	client := command.NewOnlinePvpClient(conn)

	ps := newPvpGRPCStream(client)

	for {
		if !ps.isSuccessConnected {
			ps.openStream()
		}
		time.Sleep(5 * time.Second)
	}
}

func (ps *PvpGRPCStream) openStream() {
	stream, err := ps.client.OnlinePvpStream(context.Background())
	if err != nil {
		return
	}
	ps.isSuccessConnected = true
	//fmt.Println("开启grpc流成功 online-->pvp")
	logger.Info("开启grpc流成功 online-->pvp")
	defer func() {
		ps.isSuccessConnected = false
		stream.CloseSend()
		logger.Info("grpc流断开连接 online-->pvp")
	}()
	go func() {
		for {
			select {
			case msg := <-ps.msgPvpChannel:
				if err := stream.Send(msg); err != nil {
					return
				}
			}
		}
	}()
	for {
		out, err := stream.Recv()
		if err != nil || err == io.EOF {
			return
		}
		ps.onMessage(out.PlayerId, out.Data)
	}
}

func (ps *PvpGRPCStream) onMessage(playerId uint64, data []byte) {
	WorldGetMe().sendFromOtherServerMsgChan(playerId, data)
}
