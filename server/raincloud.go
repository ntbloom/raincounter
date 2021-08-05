package main

import (
	_ "github.com/ntbloom/raincounter/common/mqtt"
	"github.com/ntbloom/raincounter/server/database"
	"github.com/ntbloom/raincounter/server/web"
)

func main() {
	database.NewDatabase("raincloud", "not-a-real-url")
	//mqtt.Start()
	web.Serve()
}
