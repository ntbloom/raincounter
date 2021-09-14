// Package mqtt wraps the Eclipse Paho code for handling paho messaging
package mqtt

import (
	"fmt"
	"time"

	"github.com/ntbloom/raincounter/pkg/config/configkey"

	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"

	paho "github.com/eclipse/paho.mqtt.golang"
)

const localhost = "127.0.0.1"
const localport = 1883

// BrokerConfig configures the mqtt connection
type BrokerConfig struct {
	broker            string
	port              int
	caCert            string
	clientCert        string
	clientKey         string
	connectionTimeout time.Duration
}

// newBrokerConfig get mqtt configuration details from viper directly
func newBrokerConfig() *BrokerConfig {
	return &BrokerConfig{
		broker:            viper.GetString(configkey.MQTTBrokerIP),
		port:              viper.GetInt(configkey.MQTTBrokerPort),
		caCert:            viper.GetString(configkey.MQTTCaCert),
		clientCert:        viper.GetString(configkey.MQTTClientCert),
		clientKey:         viper.GetString(configkey.MQTTClientKey),
		connectionTimeout: viper.GetDuration(configkey.MQTTConnectionTimeout),
	}
}

// NewConnection creates a new MQTT connection or error
func NewConnection() (paho.Client, error) {
	options := paho.NewClientOptions()
	config := newBrokerConfig()

	// add broker, authenticate if necessary
	var scheme string
	useTLS := viper.GetBool(configkey.MQTTUseTLS)
	if useTLS {
		scheme = "ssl"
	} else {
		scheme = "mqtt"
	}
	switch scheme {
	case "ssl":
		logrus.Debug("using TLS to connect")
		// configure tls
		tlsConfig, err := configureTLSConfig(config.caCert, config.clientCert, config.clientKey)
		if err != nil {
			return nil, err
		}
		options.SetTLSConfig(tlsConfig)
	case "mqtt":
		logrus.Warning("Connecting to MQTT broker on localhost:1883 without encryption, for testing only")
		config.broker = localhost
		config.port = localport
	default:
		panic(fmt.Sprintf("unsupported mqtt scheme: %s", scheme))
	}

	server := fmt.Sprintf("%s://%s:%d", scheme, config.broker, config.port)
	logrus.Debugf("opening MQTT connection at %s", server)
	options.AddBroker(server)

	// miscellaneous options
	options.SetConnectTimeout(config.connectionTimeout)

	client := paho.NewClient(options)
	return client,
		nil
}
