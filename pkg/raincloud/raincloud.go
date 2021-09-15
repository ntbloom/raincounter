package raincloud

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/ntbloom/raincounter/pkg/raincloud/receiver"
)

// Receive runs the main receiver loop
func Receive() {
	recv, err := receiver.NewReceiver()
	if err != nil {
		panic(err)
	}
	defer recv.Stop()
	go recv.Start()
	terminalSignals := make(chan os.Signal, 1)
	signal.Notify(terminalSignals, syscall.SIGINT, syscall.SIGTERM)

	sig := <-terminalSignals
	logrus.Infof("program received %s signal, exiting", sig)
	logrus.Info("Done!")
}

// Serve serves the web server
func Serve() {
	logrus.Error("implement me!")
	os.Exit(-1)
}
