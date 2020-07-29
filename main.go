package main

import (
	"landlords/helper/uuid"
	"landlords/initcards"
	"landlords/log"
	"landlords/manager"
	"landlords/websocket"
)

func main() {
	log.InitLog()

	uuid.InitUUID()

	//redis.InitRedis()

	initcards.InitNewCards()

	manager.InitRoomManager()

	manager.InitPvpPoolManager()

	//agentservice.AgentRun()
	websocket.StartWebSocket("")
}
