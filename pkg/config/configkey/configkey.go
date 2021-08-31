// Package configkey maps yaml configs to a constant
package configkey

const (
	Loglevel = "log.level"

	USBPacketLengthMax   = "usb.packet.length.max"
	USBConnectionPort    = "usb.connection.port"
	USBConnectionTimeout = "usb.connection.timeout"

	MQTTBrokerIP          = "mqtt.broker.ip"
	MQTTBrokerPort        = "mqtt.broker.port"
	MQTTCaCert            = "mqtt.certs.ca"
	MQTTClientCert        = "mqtt.certs.client"
	MQTTClientKey         = "mqtt.certs.key"
	MQTTConnectionTimeout = "mqtt.connection.timeout"
	MQTTQuiescence        = "mqtt.connection.quiescence"
	MQTTQos               = "mqtt.qos"

	SensorRainMm = "sensor.mm"

	DatabaseLocalFile          = "database.local.file"
	DatabaseRemoteName         = "database.remote.name"
	DatabasePostgresqlPassword = "database.remote.password"

	MessengerStatusInterval = "messenger.status.interval"

	MainLoopDuration = "main.loop.duration"
)

// random constants, not tied to config
const (
	DegreesF          = string("\u00B0F")
	DegreesC          = string("\u00B0C")
	Kill              = 1
	SendStatusMessage = 2
)
