package main

import (
	command "core/command/pb"
	"core/component/logger"
)

//开始游戏
func (player *Player) StartGameHandler(msgId uint16, data []byte) {
	msgRecv := &command.CSStartGame_Online{}
	err := player.innerRouter.UnMarshal(data, msgRecv)
	if err != nil {
		logger.Errorf("recv data:%v msg:%v", data, msgRecv)
		return
	}
	logger.Info(msgRecv.RoomId)
	player.sendMSg(command.Command_CSStartGame, msgRecv)
	//player.addPlayer2PvpPool(int(msgRecv.RoomId), player.User.Id, player.User.Name)
}
