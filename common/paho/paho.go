// Package paho wraps the Eclipse Paho code for handling mqtt messaging
package paho

import (
	"fmt"
	"time"

	"github.com/ntbloom/rainbase/pkg/config/configkey"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MqttConfig configures the paho connection
type MqttConfig struct {
	scheme            string
	broker            string
	port              int
	caCert            string
	clientCert        string
	clientKey         string
	connectionTimeout time.Duration
}

// GetConfigFromViper get paho configuration details from viper directly
func GetConfigFromViper() *MqttConfig {
	// look for certs locally first
	return &MqttConfig{
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
func NewConnection(config *MqttConfig) (mqtt.Client, error) {
	options := mqtt.NewClientOptions()

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

	client := mqtt.NewClient(options)
	return client,
		nil
}
