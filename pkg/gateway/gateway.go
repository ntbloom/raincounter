package gateway

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt2 "github.com/ntbloom/raincounter/pkg/common/mqtt"

	config2 "github.com/ntbloom/raincounter/pkg/config"
	configkey2 "github.com/ntbloom/raincounter/pkg/config/configkey"

	database2 "github.com/ntbloom/raincounter/pkg/gateway/database"
	messenger2 "github.com/ntbloom/raincounter/pkg/gateway/messenger"
	serial2 "github.com/ntbloom/raincounter/pkg/gateway/serial"

	paho "github.com/eclipse/paho.mqtt.golang"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// connect to mqtt
func connectToMQTT() paho.Client {
	client, err := mqtt2.NewConnection(mqtt2.NewBrokerConfig())
	if err != nil {
		panic(err)
	}
	return client
}

// connect to the sqlite database
func connectToDatabase() *database2.DBConnector {
	db, err := database2.NewSqliteDBConnector(viper.GetString(configkey2.DatabaseLocalDevFile), true)
	if err != nil {
		panic(err)
	}
	return db
}

// get a serial connection
func connectSerialPort(msgr *messenger2.Messenger) *serial2.Serial {
	conn, err := serial2.NewConnection(
		viper.GetString(configkey2.USBConnectionPort),
		viper.GetInt(configkey2.USBPacketLengthMax),
		viper.GetDuration(configkey2.USBConnectionTimeout),
		msgr,
	)
	if err != nil {
		panic(err)
	}
	return conn
}

// run launches program for seconds or indefinitely if duration is negative
func run() {
	client := connectToMQTT()
	db := connectToDatabase()
	msgr := messenger2.NewMessenger(client, db)
	conn := connectSerialPort(msgr)

	// start the listening threads
	go msgr.Start()
	go conn.Start()

	// start a timer if needed
	var loopTimer *time.Timer
	var timerChan <-chan time.Time
	duration := viper.GetDuration(configkey2.MainLoopDuration)
	if duration.Seconds() > 0 {
		loopTimer = time.NewTimer(viper.GetDuration(configkey2.MainLoopDuration))
		timerChan = loopTimer.C
	}

	// look out for terminal input
	terminalSignals := make(chan os.Signal, 1)
	signal.Notify(terminalSignals, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-terminalSignals:
			logrus.Infof("program received %s signal, exiting", sig)
			stopProgram(msgr, conn, loopTimer)
		case <-timerChan:
			logrus.Infof("program exiting after %s", duration)
			stopProgram(msgr, conn, loopTimer)
		}
	}
}

func stopProgram(msgr *messenger2.Messenger, conn *serial2.Serial, timer *time.Timer) {
	if timer != nil {
		timer.Stop()
	}
	msgr.Stop()
	conn.Stop()

	time.Sleep(time.Second * 1)
	logrus.Info("Done!")
	os.Exit(0)
}

func Start() {
	// read config from the config file
	config2.Configure()

	// run the main listening loop
	run()
}
