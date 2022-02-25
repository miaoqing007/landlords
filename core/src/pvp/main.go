package main

import (
	"core/component/logger"
)

func main() {
	logger.SetLogFile("../log/pvp"+"_"+"1", "pvp")
	logger.SetLogLevel("DEBUG")

	runOnlinePvpGRPC()
}
