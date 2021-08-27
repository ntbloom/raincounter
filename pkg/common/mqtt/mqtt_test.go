package mqtt_test

import (
	"testing"

	"github.com/ntbloom/raincounter/pkg/common/mqtt"

	"github.com/ntbloom/raincounter/pkg/config"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

// reusable mqtt function
func pahoFixture(t *testing.T) paho.Client {
	config.Configure()
	pahoConfig := mqtt.NewBrokerConfig(true)
	client, err := mqtt.NewConnection(pahoConfig)
	if err != nil {
		t.Fail()
	}
	return client
}

// Can we connect with the remote server (requires server to be working)
func TestMQTTConnection(t *testing.T) {
	client := pahoFixture(t)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		t.Fail()
	}
	defer client.Disconnect(1000)
	if !client.IsConnected() {
		logrus.Error("failed to connect")
		t.Fail()
	}
	client.Publish("hello", 0, false, "world")
}
