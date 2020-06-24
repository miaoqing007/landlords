package main

import (
	"landlords/agentservice"
	"landlords/helper/uuid"
	"landlords/initcards"
	"landlords/log"
	"landlords/manager"
)

func main() {
	log.InitLog()

	uuid.InitUUID()

	//redis.InitRedis()

	initcards.InitNewCards()

	manager.InitRoomManager()

	manager.InitPvpPoolManager()

	agentservice.AgentRun()
}
