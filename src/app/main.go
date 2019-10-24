package main

import (
	"flag"
)

var host = flag.String("host", "", "host")
var port = flag.String("port", "9999", "port")

func main() {
	agentRun()

}
