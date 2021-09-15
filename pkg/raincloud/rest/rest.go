package rest

import "github.com/sirupsen/logrus"

// Serve runs rest server
func Serve() {
	logrus.Info("this is where we serve a read-only REST endpoint for the website from the postgresql")
}
