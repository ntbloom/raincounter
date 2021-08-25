// package web runs a web REST API
package web

import "github.com/sirupsen/logrus"

// Serve runs web server
func Serve() {
	logrus.Info("this is where we serve a read-only REST endpoint for the website from the postgresql")
}
