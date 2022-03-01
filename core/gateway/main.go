package main

import (
	"core/component/logger"
	"flag"
)

func main() {
	flag.Parse()
	logger.SetLogFile("../log/gateway"+"_"+"1", "gateway")
	logger.SetLogLevel("DEBUG")

	runTcpAndGRPC()
}
