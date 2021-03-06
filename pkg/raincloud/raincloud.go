package raincloud

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ntbloom/raincounter/pkg/raincloud/frontend"

	"github.com/sirupsen/logrus"

	"github.com/ntbloom/raincounter/pkg/raincloud/receiver"
)

// wait for sigint or sigterm before quitting
func waitForSignal() {
	terminalSignals := make(chan os.Signal, 1)
	signal.Notify(terminalSignals, syscall.SIGINT, syscall.SIGTERM)

	sig := <-terminalSignals
	logrus.Infof("program received %s signal, exiting", sig)
	logrus.Info("Done!")
}

// Receive runs the main receiver loop
func Receive() {
	recv, err := receiver.NewReceiver()
	if err != nil {
		panic(err)
	}
	defer recv.Stop()
	go recv.Start()

	waitForSignal()
}

// Serve serves the web server
func Serve() {
	server, err := frontend.NewHTMLServer()
	if err != nil {
		panic(err)
	}
	go server.Run()
	defer server.Stop()

	waitForSignal()
}
