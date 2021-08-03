package main

import (
	_ "github.com/ntbloom/raincounter/common/paho"
	"github.com/ntbloom/raincounter/server/database"
	"github.com/ntbloom/raincounter/server/web"
)

func main() {
	database.NewDatabase("raincloud", "not-a-real-url")
	//mqtt.Loop()
	web.Serve()
}
