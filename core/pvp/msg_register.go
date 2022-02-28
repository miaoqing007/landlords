package main

import command "core/command/pb"

func (player *Player) registerMsgHandler() {
	player.router.RegisterOnlinePvp(uint16(command.Command_CSJoinPvpPool), player.joinPvpPoolHandler)
	player.router.RegisterOnlinePvp(uint16(command.Command_CSExitPvpPool), player.exitPvpPoolHandler)
}
