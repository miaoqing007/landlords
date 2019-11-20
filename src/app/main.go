package main

import (
	"app/initcards"
	"app/log"
	"app/manager"
	"app/redis"
)

func main() {
	log.InitLog()

	redis.InitRedis()

	initcards.InitNewCards()

	manager.InitRoomManager()

	manager.InitPvpPoolManager()

	agentRun()
}
