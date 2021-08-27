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
func NewBrokerConfig(withAuth bool) *BrokerConfig {
	// look for certs locally first
	return &BrokerConfig{
		broker:            viper.GetString(configkey.MQTTBrokerIP),
		port:              viper.GetInt(configkey.MQTTBrokerPort),
		caCert:            viper.GetString(configkey.MQTTCaCert),
		clientCert:        viper.GetString(configkey.MQTTClientCert),
		clientKey:         viper.GetString(configkey.MQTTClientKey),
		connectionTimeout: viper.GetDuration(configkey.MQTTConnectionTimeout),
		auth:              withAuth,
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

// LocalDevConnection gets you attached to a local docker container without auth
func LocalDevConnection(hostname string, port int) paho.Client {
	options := paho.NewClientOptions()
	server := fmt.Sprintf("mqtt://%s:%d", hostname, port)
	options.AddBroker(server)

	client := paho.NewClient(options)
	return client
}
