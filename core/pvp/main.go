package main

import (
	"core/component/logger"
	"flag"
)

func main() {
	flag.Parse()
	logger.SetLogFile("../log/pvp"+"_"+"1", "pvp")
	logger.SetLogLevel("DEBUG")

	runOnlinePvpGRPC()
}
