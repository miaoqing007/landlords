package landlords

import (
	"landlords/helper/uuid"
	"landlords/initcards"
	"landlords/log"
	"landlords/manager"
	"landlords/redis"
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
