package main

import "core/component/logger"

func main() {
	logger.SetLogFile("../log/gateway"+"_"+"1", "gateway")
	logger.SetLogLevel("DEBUG")

	runTcpAndGRPC()

}
