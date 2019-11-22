package main

import (
	"app/helper/uuid"
	"app/initcards"
	"app/log"
	"app/manager"
	"app/redis"
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
