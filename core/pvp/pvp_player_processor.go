package main

import (
	command "core/command/pb"
	"core/component/logger"
)

func (player *Player) joinPvpPoolHandler(msgID uint16, data []byte) {
	msgRecv := &command.CSJoinPvpPool_Pvp{}
	err := player.router.UnMarshal(data, msgRecv)
	if err != nil {
		logger.Errorf("recv data:%v msg:%v", data, msgRecv)
		return
	}
	logger.Info(msgRecv.PlayerId)
	//player.sendMsgToOnlineClient()
}

func (player *Player) exitPvpPoolHandler(msgID uint16, data []byte) {
	msgRecv := &command.CSJoinPvpPool_Pvp{}
	err := player.router.UnMarshal(data, msgRecv)
	if err != nil {
		logger.Errorf("recv data:%v msg:%v", data, msgRecv)
		return
	}
	//player.sendMsgToOnlineClient()
}
