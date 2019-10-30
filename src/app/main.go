package main

import (
	"app/initcards"
	"app/log"
	"app/manager"
)

func main() {
	log.InitLog()
	initcards.InitNewCards()
	manager.InitRoomManager()
	agentRun()
}
