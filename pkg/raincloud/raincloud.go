package raincloud

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ntbloom/raincounter/pkg/raincloud/api"

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
	rest, err := api.NewRestServer()
	if err != nil {
		logrus.Fatal("unable to run the rest API")
		panic(err)
	}
	var duration time.Duration = 10
	logrus.Infof("In debug mode, running for %d seconds", duration)
	go rest.Run()
	time.Sleep(time.Second * duration)
	rest.Stop()
}
