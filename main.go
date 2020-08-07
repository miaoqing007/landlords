package main

import (
	"landlords/config"
	"landlords/initcards"
	"landlords/log"
	"landlords/manager"
	"landlords/redis"
	"landlords/signal"
	"landlords/websocket"
)

func main() {
	log.InitLog()

	config.InitConfig()

	redis.InitRedis()

	signal.InitSignal()

	initcards.InitNewCards()

	manager.InitRoomManager()

	manager.InitPvpPoolManager()

	websocket.Run()
}
