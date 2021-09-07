package mqtt

import (
	"time"

	"github.com/ntbloom/raincounter/pkg/gateway/messenger"
)

// Examples of sample messages for use in testing, etc.

type SampleMessage struct {
	topic string
	msg   map[string]interface{}
}

var timestamp = time.Now()

var SampleRain = SampleMessage{
	topic: RainTopic,
	msg:   map[string]interface{}{"Millimeters": "0.2794", "Timestamp": timestamp},
}

var SampleTemp = SampleMessage{
	topic: TemperatureTopic,
	msg:   map[string]interface{}{"Tempc": 23, "Timestamp": timestamp},
}

var SampleSensorPause = SampleMessage{
	topic: SensorEvent,
	msg:   map[string]interface{}{"Status": messenger.SensorPause, "Timestamp": timestamp},
}

var SampleSensorUnpause = SampleMessage{
	topic: SensorEvent,
	msg:   map[string]interface{}{"Status": messenger.SensorUnpause, "Timestamp": timestamp},
}

var SampleSensorSoftReset = SampleMessage{
	topic: SensorEvent,
	msg:   map[string]interface{}{"Status": messenger.SensorSoftReset, "Timestamp": timestamp},
}

var SampleSensorHardReset = SampleMessage{
	topic: SensorEvent,
	msg:   map[string]interface{}{"Status": messenger.SensorHardReset, "Timestamp": timestamp},
}

var SampleSensorStatus = SampleMessage{
	topic: SensorStatus,
	msg:   map[string]interface{}{"OK": true, "Timestamp": timestamp},
}

var SampleGatewayStatus = SampleMessage{
	topic: GatewayStatus,
	msg:   map[string]interface{}{"OK": true, "Timestamp": timestamp},
}
