package main

import (
	"core/online/log"
	"core/online/redis"
)

func main() {
	log.InitLog()
	redis.InitRedis()
	runGatewayOnlineGRPC()
}

//
//import (
//	"agentservice"
//	"landlords/config"
//	"landlords/initcards"
//	"landlords/log"
//	"landlords/manager"
//	"landlords/redis"
//	"landlords/signal"
//)
//
//func main() {
//	log.InitLog()
//
//	config.InitConfig()
//
//	redis.InitRedis()
//
//	signal.InitSignal()
//
//	initcards.InitNewCards()
//
//	manager.InitRoomManager()
//
//	manager.InitPvpPoolManager()
//
//	//websocket.Run()
//	agentservice.AgentRun()
//}
