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

// BrokerConfig configures the mqtt connection
type BrokerConfig struct {
	scheme            string
	broker            string
	port              int
	caCert            string
	clientCert        string
	clientKey         string
	connectionTimeout time.Duration
}

// NewBrokerConfig get mqtt configuration details from viper directly
func NewBrokerConfig() *BrokerConfig {
	// look for certs locally first
	return &BrokerConfig{
		scheme:            viper.GetString(configkey.MQTTScheme),
		broker:            viper.GetString(configkey.MQTTBrokerIP),
		port:              viper.GetInt(configkey.MQTTBrokerPort),
		caCert:            viper.GetString(configkey.MQTTCaCert),
		clientCert:        viper.GetString(configkey.MQTTClientCert),
		clientKey:         viper.GetString(configkey.MQTTClientKey),
		connectionTimeout: viper.GetDuration(configkey.MQTTConnectionTimeout),
	}
}

// NewConnection creates a new MQTT connection or error
func NewConnection(config *BrokerConfig) (paho.Client, error) {
	options := paho.NewClientOptions()

	// add broker
	server := fmt.Sprintf("%s://%s:%d", config.scheme, config.broker, config.port)
	logrus.Debugf("opening MQTT connection at %s", server)
	options.AddBroker(server)

	// configure tls
	tlsConfig, err := configureTLSConfig(config.caCert, config.clientCert, config.clientKey)
	if err != nil {
		return nil, err
	}
	options.SetTLSConfig(tlsConfig)

	// miscellaneous options
	options.SetConnectTimeout(config.connectionTimeout)

	client := paho.NewClient(options)
	return client,
		nil
}
