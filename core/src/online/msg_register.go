package main

import (
	"landlords/client_handler"
	command "landlords/command/pb"
)

func (player *Player) RegisterPlayerMsg() {
	player.router.Register(uint16(command.Command_CSStartGame), client_handler.StartGameHandler)
}
