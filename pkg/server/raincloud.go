package main

import (
	_ "github.com/ntbloom/raincounter/pkg/common/mqtt"
	database2 "github.com/ntbloom/raincounter/pkg/server/postgresql"
	web2 "github.com/ntbloom/raincounter/pkg/server/web"
)

func main() {
	database2.NewDatabase("raincloud", "not-a-real-url")
	//mqtt.Start()
	web2.Serve()
}
