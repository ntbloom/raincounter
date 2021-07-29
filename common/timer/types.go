package timer

import "time"

// Define some pre-canned timer types

type channelUint8Timer struct {
	channel chan uint8
	value   uint8
}

func (ch *channelUint8Timer) DoAction() {
	ch.channel <- ch.value
}

// NewChannelUint8Timer creates a timer that sends `value` to `channel` every `duration` period
func NewChannelUint8Timer(interval, frequency time.Duration, channel chan uint8, value uint8) *Timer {
	ch := &channelUint8Timer{channel, value}
	t := NewTimer(interval, frequency, ch)
	return t
}
