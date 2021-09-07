package mqtt

// Examples of sample messages for use in testing, etc.

type SampleMessage struct {
	topic string
	msg   map[string]interface{}
}

const timestamp = "2021-09-06T21:57:32.779567444-04:00"

var SampleRain = SampleMessage{
	topic: RainTopic,
	msg:   map[string]interface{}{"Millimeters": "0.2794", "Timestamp": timestamp},
}

var SampleTemp = SampleMessage{
	topic: TemperatureTopic,
	msg:   map[string]interface{}{"Tempc": 23, "Timestamp": timestamp},
}

var SampleSensorPause = SampleMessage{
	topic: SensorEventTopic,
	msg:   map[string]interface{}{"Status": SensorPauseEvent, "Timestamp": timestamp},
}

var SampleSensorUnpause = SampleMessage{
	topic: SensorEventTopic,
	msg:   map[string]interface{}{"Status": SensorUnpauseEvent, "Timestamp": timestamp},
}

var SampleSensorSoftReset = SampleMessage{
	topic: SensorEventTopic,
	msg:   map[string]interface{}{"Status": SensorSoftResetEvent, "Timestamp": timestamp},
}

var SampleSensorHardReset = SampleMessage{
	topic: SensorEventTopic,
	msg:   map[string]interface{}{"Status": SensorHardResetEvent, "Timestamp": timestamp},
}

var SampleSensorStatus = SampleMessage{
	topic: SensorStatusTopic,
	msg:   map[string]interface{}{"OK": true, "Timestamp": timestamp},
}

var SampleGatewayStatus = SampleMessage{
	topic: GatewayStatusTopic,
	msg:   map[string]interface{}{"OK": true, "Timestamp": timestamp},
}
