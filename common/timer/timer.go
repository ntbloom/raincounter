// package timer executes functions on a schedule

package timer

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Action is a base interface for implementing callbacks that occur after a given
// duration. Calls to `DoAction()` occur when timers run out.
type Action interface {
	DoAction() // run a parameterless function
}

// Timer triggers actions at regular intervals
type Timer struct {
	start     time.Time     // when the timer starts
	interval  time.Duration // when to trigger something
	frequency time.Duration // how often to check the clock
	action    Action        // function to call when timer is up
	Kill      chan bool     // send bool to channel to finish
}

// NewTimer returns a pointer to a Timer struct
func NewTimer(interval, frequency time.Duration, action Action) *Timer {
	if frequency > interval {
		logrus.Errorf("%s > %s", frequency, interval)
		panic("frequency must be less than interval")
	}
	finish := make(chan bool)
	return &Timer{
		start:     time.Now(),
		interval:  interval,
		frequency: frequency,
		action:    action,
		Kill:      finish,
	}
}

// Loop runs an infinite loop, triggering an action. stops when receives message on Kill channel
func (t *Timer) Loop() {
	trigger := make(chan bool)
	stopChecking := make(chan bool)
	if t.interval > 0 {
		go t.checkTimer(trigger, stopChecking)
	}

	for {
		select {
		case <-t.Kill:
			return
		case <-trigger:
			go t.action.DoAction()
		}
	}
}

// checkTimer infinitely checks clock every `wait` duration, sends message when trigger is up
func (t *Timer) checkTimer(trigger, stopChecking chan bool) {
	for {
		if time.Since(t.start) > t.interval {
			t.start = time.Now()
			trigger <- true
		}
		select {
		case <-stopChecking:
			return
		default:
			time.Sleep(t.frequency)
		}
	}
}
