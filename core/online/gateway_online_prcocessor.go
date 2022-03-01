package main

import (
	command "core/command/pb"
	"core/component/logger"
)

//玩家进入online
func (os *OnlineStreamInfo) clientInOnlineHandler(msgId uint16, data []byte) {
	msgRecv := &command.ClientInOnline_Online{}
	err := os.router.UnMarshal(data, msgRecv)
	if err != nil {
		logger.Errorf("recv data:%v msg:%v", data, msgRecv)
		return
	}
	WorldGetMe().addPlayer(msgRecv.PlayerId, msgRecv.ClientAddr, os.gatewayRemoteAddr)
}

//玩家退出online
func (os *OnlineStreamInfo) clientOutOnlineHandler(msgId uint16, data []byte) {
	msgRecv := &command.ClientOutOnline_Online{}
	err := os.router.UnMarshal(data, msgRecv)
	if err != nil {
		logger.Errorf("recv data:%v msg:%v", data, msgRecv)
		return
	}
	if player := WorldGetMe().getPlayer(msgRecv.PlayerId); player != nil {
		player.addCloseChannel()
	}
}
