package main

import (
	"landlords/app/helper/uuid"
	"landlords/app/initcards"
	"landlords/app/log"
	"landlords/app/manager"
	"landlords/app/redis"
)

func main() {
	log.InitLog()

	uuid.InitUUID()

	redis.InitRedis()

	initcards.InitNewCards()

	manager.InitRoomManager()

	manager.InitPvpPoolManager()

	agentRun()
}
