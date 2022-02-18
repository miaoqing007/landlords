package main

import (
	command "landlords/command/pb"
)

//开始游戏
func (player *Player) StartGameHandler(msgId uint16, data []byte) {
	msgRecv := &command.CSStartGame{}
	err := player.router.UnMarshal(data, msgRecv)
	if err != nil {
		return
	}
	player.addPlayer2PvpPool(int(msgRecv.RoomId), player.User.Id, player.User.Name)
}
