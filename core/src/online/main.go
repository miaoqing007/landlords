package main

func main() {
	runGRPCDial("127.0.0.1:9999")
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
