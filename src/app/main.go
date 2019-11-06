package main

import (
	"app/client_handler"
	"app/initcards"
	"app/log"
	"app/manager"
	"app/redisgo"
)

func main() {
	log.InitLog()

	redisgo.InitRedis()

	client_handler.InitHandle()

	initcards.InitNewCards()

	manager.InitRoomManager()

	agentRun()
}
