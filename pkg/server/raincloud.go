package raincloud

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/ntbloom/raincounter/pkg/server/receiver"
)

func Start() {
	recv, err := receiver.NewReceiver()
	if err != nil {
		panic(err)
	}
	defer recv.Close()
	go recv.Start()
	terminalSignals := make(chan os.Signal, 1)
	signal.Notify(terminalSignals, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case sig := <-terminalSignals:
			logrus.Infof("program received %s signal, exiting", sig)
			recv.Stop()
		default:
			continue
		}
	}
}
