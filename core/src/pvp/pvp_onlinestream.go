package main

import (
	command "core/command/pb"
	"google.golang.org/grpc/peer"
	"net"
)

type OnlineStreamInfo struct {
	onlineStream         command.OnlinePvp_OnlinePvpStreamServer
	onlineGRPCRemoteAddr string
	msgToOnline          chan *command.Pvp2OnlineInfo
}

func newOnlineStreamInfo(streamServer command.OnlinePvp_OnlinePvpStreamServer) *OnlineStreamInfo {
	os := &OnlineStreamInfo{
		onlineStream: streamServer,
		msgToOnline:  make(chan *command.Pvp2OnlineInfo, 1024),
	}
	os.onlineGRPCRemoteAddr = os.getOnlineGRPCRemoteAddr()
	return os
}

func (os *OnlineStreamInfo) addMsgToOnine(playerId uint64, data []byte) {
	os.msgToOnline <- &command.Pvp2OnlineInfo{Data: data, PlayerId: playerId}
}

func (os *OnlineStreamInfo) getOnlineGRPCRemoteAddr() string {
	pr, ok := peer.FromContext(os.onlineStream.Context())
	if !ok {
		return ""
	}
	if pr.Addr == net.Addr(nil) {
		return ""
	}
	return pr.Addr.String()
}
