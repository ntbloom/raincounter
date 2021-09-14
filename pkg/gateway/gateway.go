package gateway

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ntbloom/raincounter/pkg/gateway/localdb"

	"github.com/ntbloom/raincounter/pkg/common/mqtt"

	"github.com/ntbloom/raincounter/pkg/config"
	"github.com/ntbloom/raincounter/pkg/config/configkey"

	"github.com/ntbloom/raincounter/pkg/gateway/messenger"
	"github.com/ntbloom/raincounter/pkg/gateway/serial"

	paho "github.com/eclipse/paho.mqtt.golang"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// connect to mqtt
func connectToMQTT() paho.Client {
	client, err := mqtt.NewConnection(mqtt.newBrokerConfig())
	if err != nil {
		panic(err)
	}
	return client
}

// connect to the localdb postgresql
func connectToDatabase() *localdb.LocalDB {
	db, err := localdb.NewLocalDB(viper.GetString(configkey.DatabaseLocalFile), true)
	if err != nil {
		panic(err)
	}
	return db
}

// get a serial connection
func connectSerialPort(msgr *messenger.Messenger) *serial.Serial {
	conn, err := serial.NewConnection(
		viper.GetString(configkey.USBConnectionPort),
		viper.GetInt(configkey.USBPacketLengthMax),
		viper.GetDuration(configkey.USBConnectionTimeout),
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
	msgr := messenger.NewMessenger(client, db)
	conn := connectSerialPort(msgr)

	// start the listening threads
	go msgr.Start()
	go conn.Start()

	// start a timer if needed
	var loopTimer *time.Timer
	var timerChan <-chan time.Time
	duration := viper.GetDuration(configkey.MainLoopDuration)
	if duration.Seconds() > 0 {
		loopTimer = time.NewTimer(viper.GetDuration(configkey.MainLoopDuration))
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

func stopProgram(msgr *messenger.Messenger, conn *serial.Serial, timer *time.Timer) {
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
	config.Configure()

	// run the main listening loop
	run()
}
