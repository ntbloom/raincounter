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
	MQTTUseTLS            = "mqtt.tls"
	MQTTClientCert        = "mqtt.certs.client"
	MQTTClientKey         = "mqtt.certs.key"
	MQTTConnectionTimeout = "mqtt.connection.timeout"
	MQTTQuiescence        = "mqtt.connection.quiescence"
	MQTTQos               = "mqtt.qos"

	SensorRainMm        = "sensor.mm"
	AssetStatusDuration = "asset.status.duration"

	DatabaseLocalFile = "database.local.file"

	PGDatabaseName        = "database.remote.name"
	PGPassword            = "database.remote.password"
	PGConnectionTimeout   = "database.remote.connection.timeout"
	PGConnectionRetryWait = "database.remote.connection.retry.wait"

	MessengerStatusInterval = "messenger.status.interval"

	MainLoopDuration = "main.loop.duration"

	RestIP      = "rest.ip.address"
	RestPort    = "rest.ip.port"
	RestScheme  = "rest.scheme"
	RestVersion = "rest.version"

	WebEntrypoint    = "web.entrypoint"
	WebDirectory     = "web.directory"
	WebServerAddress = ":8080"
)

// random constants, not tied to config
const (
	DegreesF          = string("\u00B0F")
	DegreesC          = string("\u00B0C")
	Kill              = 1
	SendStatusMessage = 2
	TimestampFormat   = time.RFC3339
	PrettyTimeFormat  = "Mon Jan 2, 2006 15:04 EST"
	SensorStatus      = 1
	GatewayStatus     = 2
	IntErrVal         = -99
	FloatErrVal       = -999.0
)
