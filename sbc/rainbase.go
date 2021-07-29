package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/ntbloom/rainbase/pkg/config"
	"github.com/ntbloom/rainbase/pkg/config/configkey"
	"github.com/ntbloom/rainbase/pkg/database"
	"github.com/ntbloom/rainbase/pkg/messenger"
	"github.com/ntbloom/rainbase/pkg/paho"
	"github.com/ntbloom/rainbase/pkg/serial"
	"github.com/ntbloom/rainbase/pkg/timer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// connect to paho
func connectToMQTT() mqtt.Client {
	client, err := paho.NewConnection(paho.GetConfigFromViper())
	if err != nil {
		panic(err)
	}
	return client
}

// connect to the sqlite database
func connectToDatabase() *database.DBConnector {
	db, err := database.NewSqliteDBConnector(viper.GetString(configkey.DatabaseLocalDevFile), true)
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

// set up a way to kill the process, either through a timer or os.signal
func stopLoop(channels []chan uint8) {
	for _, channel := range channels {
		channel <- configkey.SerialClosed
	}
}

// struct for killing with timer
type kill struct{ channels []chan uint8 }

func (k *kill) DoAction() {
	stopLoop(k.channels)
	logrus.Info("main loop killed by timer, exiting program")
	os.Exit(0)
}

func startKillTimer(duration, frequency time.Duration, killChannels []chan uint8) *timer.Timer {
	k := kill{killChannels}
	t := timer.NewTimer(duration, frequency, &k)
	return t
}

// run main listening loop for number of seconds or indefinitely if duration is negative
func listen(duration, frequency time.Duration) {
	client := connectToMQTT()
	db := connectToDatabase()
	msgr := messenger.NewMessenger(client, db)
	conn := connectSerialPort(msgr)

	// start the listening threads
	go msgr.Listen()
	go conn.GetTLV()

	// start a timer
	killChannels := []chan uint8{conn.State, msgr.State}
	t := startKillTimer(duration, frequency, killChannels)
	go t.Loop()

	// kill process with sigint regardless of whether duration is negative
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		logrus.Infof("program received %s signal, exiting", sig)
		stopLoop(killChannels)
		t.Kill <- true
	}()
	<-t.Kill
}

func main() {
	// read config from the config file
	config.Configure()

	// run the main listening loop
	duration := viper.GetDuration(configkey.MainLoopDuration)
	frequency := viper.GetDuration(configkey.MainLoopFrequency)
	listen(duration, frequency)
}
