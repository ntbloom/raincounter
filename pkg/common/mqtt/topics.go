package mqtt

// mqtt topics for events, measurements, status
const (
	GatewayStatusTopic = "status/gateway"
	SensorStatusTopic  = "status/sensor"
	TemperatureTopic   = "measurement/temperature"
	RainTopic          = "measurement/rain"
	SensorEventTopic   = "sensor/event"
)

// mqtt event tags
const (
	SensorPauseEvent     = "sensorPause"
	SensorUnpauseEvent   = "sensorUnpause"
	SensorSoftResetEvent = "sensorSoftReset"
	SensorHardResetEvent = "sensorHardReset"
)
