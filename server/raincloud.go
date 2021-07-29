package main

import (
	"github.com/ntbloom/raincloud/pkg/database"
	"github.com/ntbloom/raincloud/pkg/mqtt"
	"github.com/ntbloom/raincloud/pkg/web"
)

func main() {
	database.NewDatabase("raincloud", "not-a-real-url")
	mqtt.Listen()
	web.Serve()
}
