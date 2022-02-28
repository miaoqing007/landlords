package main

import (
	"core/component/logger"
	"flag"
)

func main() {
	flag.Parse()
	logger.SetLogFile("../log/online"+"_"+"1", "online")
	logger.SetLogLevel("DEBUG")
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
