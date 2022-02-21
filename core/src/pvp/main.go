package main

import (
	//"landlords/src/config"
	//"landlords/src/online/log"
	//"landlords/src/online/redis"
	"core/config"
	"core/online/log"
	"core/online/redis"
)

func main() {
	config.InitConfig()
	log.InitLog()
	redis.InitRedis()
	runOnlinePvpGRPC()
}
