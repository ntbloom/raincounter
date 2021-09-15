package raincloud

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/ntbloom/raincounter/pkg/server/receiver"
)

// Start runs the main receiver loop
func Start() {
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
