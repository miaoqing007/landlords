package main

import (
	"app/client_handler"
	"app/initcards"
	"app/log"
	"app/manager"
	"app/model"
)

func main() {
	log.InitLog()

	model.InitRedis()

	client_handler.InitHandle()

	initcards.InitNewCards()

	manager.InitRoomManager()

	agentRun()
}
