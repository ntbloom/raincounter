// Package configkey maps yaml configs to a constant
package configkey

import "time"

const (
	Loglevel = "log.level"

	USBPacketLengthMax   = "usb.packet.length.max"
	USBConnectionPort    = "usb.connection.port"
	USBConnectionTimeout = "usb.connection.timeout"

	MQTTBrokerIP          = "mqtt.broker.ip"
	MQTTBrokerPort        = "mqtt.broker.port"
	MQTTCaCert            = "mqtt.certs.ca"
	MQTTScheme            = "mqtt.scheme"
	MQTTClientCert        = "mqtt.certs.client"
	MQTTClientKey         = "mqtt.certs.key"
	MQTTConnectionTimeout = "mqtt.connection.timeout"
	MQTTQuiescence        = "mqtt.connection.quiescence"
	MQTTQos               = "mqtt.qos"

	SensorRainMm = "sensor.mm"

	DatabaseLocalFile = "database.local.file"

	PGDatabaseName        = "database.remote.name"
	PGPassword            = "database.remote.password"
	PGConnectionTimeout   = "database.remote.connection.timeout"
	PGConnectionRetryWait = "database.remote.connection.retry.wait"

	MessengerStatusInterval = "messenger.status.interval"

	MainLoopDuration = "main.loop.duration"
)

// random constants, not tied to config
const (
	DegreesF          = string("\u00B0F")
	DegreesC          = string("\u00B0C")
	Kill              = 1
	SendStatusMessage = 2
	TimestampFormat   = time.RFC3339
	SensorStatus      = 1
	GatewayStatus     = 2
	IntErrVal         = -99
	FloatErrVal       = -999.0
)
