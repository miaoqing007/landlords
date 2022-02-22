package main

import (
	command "core/command/pb"
)

//开始游戏
func (player *Player) StartGameHandler(msgId uint16, data []byte) {
	msgRecv := &command.CSStartGameOnline{}
	err := player.innerRouter.UnMarshal(data, msgRecv)
	if err != nil {
		return
	}
	//player.addPlayer2PvpPool(int(msgRecv.RoomId), player.User.Id, player.User.Name)
}
