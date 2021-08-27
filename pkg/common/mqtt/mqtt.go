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
	broker            string
	port              int
	caCert            string
	clientCert        string
	clientKey         string
	connectionTimeout time.Duration
	auth              bool
}

// NewBrokerConfig get mqtt configuration details from viper directly
func NewBrokerConfig() *BrokerConfig {
	return &BrokerConfig{
		broker:            viper.GetString(configkey.MQTTBrokerIP),
		port:              viper.GetInt(configkey.MQTTBrokerPort),
		caCert:            viper.GetString(configkey.MQTTCaCert),
		clientCert:        viper.GetString(configkey.MQTTClientCert),
		clientKey:         viper.GetString(configkey.MQTTClientKey),
		connectionTimeout: viper.GetDuration(configkey.MQTTConnectionTimeout),
		auth:              true,
	}
}

// NewBrokerConfigNoAuth broker config with no auth, for testing only
func NewBrokerConfigNoAuth(host string, port int) *BrokerConfig {
	return &BrokerConfig{
		broker:            host,
		port:              port,
		caCert:            "/dev/null",
		clientCert:        "/dev/null",
		clientKey:         "/dev/null",
		connectionTimeout: viper.GetDuration(configkey.MQTTConnectionTimeout),
		auth:              false,
	}
}

// NewConnection creates a new MQTT connection or error
func NewConnection(config *BrokerConfig) (paho.Client, error) {
	options := paho.NewClientOptions()

	// add broker, authenticate if necessary
	scheme := "mqtt"
	if config.auth {
		scheme = "ssl"
		logrus.Debug("using TLS to connect")
		// configure tls
		tlsConfig, err := configureTLSConfig(config.caCert, config.clientCert, config.clientKey)
		if err != nil {
			return nil, err
		}
		options.SetTLSConfig(tlsConfig)
	} else {
		logrus.Debug("skipping TLS on connection")
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
