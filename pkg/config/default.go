package config

import (
	"time"

	"github.com/ntbloom/raincounter/pkg/config/configkey"
	"github.com/sirupsen/logrus"
)

var DefaultConfig = map[string]interface{}{
	configkey.Loglevel:                logrus.InfoLevel,
	configkey.USBPacketLengthMax:      7, //nolint:gomnd
	configkey.USBConnectionPort:       "/dev/ttyACM99",
	configkey.USBConnectionTimeout:    time.Second * 10,
	configkey.MQTTScheme:              "ssl",
	configkey.MQTTBrokerIP:            "127.0.0.1",
	configkey.MQTTBrokerPort:          "1883",
	configkey.MQTTCaCert:              "/etc/raincounter/ssl/client/ca.pem",
	configkey.MQTTClientCert:          "/etc/raincounter/ssl/client/client.crt",
	configkey.MQTTClientKey:           "/etc/raincounter/ssl/client/client.key",
	configkey.MQTTConnectionTimeout:   time.Second * 5,
	configkey.MQTTQuiescence:          1000,
	configkey.MQTTQos:                 1,
	configkey.SensorRainMm:            0.2794,
	configkey.DatabaseLocalFile:       "/tmp/rainbase.db",
	configkey.DatabaseRemoteFile:      "/tmp/raincloud.db",
	configkey.MessengerStatusInterval: time.Second * 10,
	configkey.MainLoopDuration:        time.Second * -10,
}
