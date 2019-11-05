package main

import (
	"app/client_handler"
	"app/initcards"
	"app/log"
	"app/manager"
)

func main() {
	log.InitLog()
	client_handler.InitHandle()
	initcards.InitNewCards()
	manager.InitRoomManager()
	agentRun()
}
