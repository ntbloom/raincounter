package raincloud

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ntbloom/raincounter/pkg/config"

	"github.com/sirupsen/logrus"

	"github.com/ntbloom/raincounter/pkg/common/mqtt"
	"github.com/ntbloom/raincounter/pkg/server/receiver"
)

func Start() {
	config.Configure()

	client, err := mqtt.NewConnection()
	if err != nil {
		panic(err)
	}
	recv, err := receiver.NewReceiver(client)
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
