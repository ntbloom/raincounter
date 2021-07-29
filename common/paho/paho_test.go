package paho_test

import (
	"testing"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"

	"github.com/ntbloom/rainbase/pkg/config"

	"github.com/ntbloom/rainbase/pkg/paho"
)

// reusable paho function
func pahoFixture(t *testing.T) mqtt.Client {
	config.Configure()
	pahoConfig := paho.GetConfigFromViper()
	client, err := paho.NewConnection(pahoConfig)
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
