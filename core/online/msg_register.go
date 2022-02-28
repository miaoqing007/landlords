package main

import (
	command "core/command/pb"
)

func (player *Player) regMsgHandler() {
	player.innerRouter.Register(uint16(command.Command_CSStartGame), player.StartGameHandler)
}
