package signal

import (
	"os"
	"os/signal"
	"syscall"
)

func InitSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGKILL, syscall.SIGINT)
	go listen(ch)
}

func listen(ch chan os.Signal) {
	for {
		select {
		case <-ch:
			os.Exit(0)
		}
	}
}
