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
	name      string        // short name for the timer
	start     time.Time     // when the timer starts
	interval  time.Duration // when to trigger something
	frequency time.Duration // how often to check the clock
	action    Action        // function to call when timer is up
	Kill      chan bool     // send bool to channel to finish
}

// NewTimer returns a pointer to a Timer struct
func NewTimer(name string, interval, frequency time.Duration, action Action) *Timer {
	logrus.Debugf("starting a new timer for %s", name)
	if interval > 0 && frequency > interval {
		logrus.Errorf("%s > %s", frequency, interval)
		panic("frequency must be less than interval for all positive interval values")
	}
	finish := make(chan bool)
	return &Timer{
		name:      name,
		start:     time.Now(),
		interval:  interval,
		frequency: frequency,
		action:    action,
		Kill:      finish,
	}
}

// Loop runs an infinite loop, triggering an action. stops when receives message on Kill channel
func (t *Timer) Loop() {
	trigger := make(chan bool, 1)
	kill := make(chan bool, 1)
	if t.interval > 0 {
		go t.checkTimer(trigger, kill)
	}

	for {
		select {
		case <-t.Kill:
			logrus.Debugf("received kill signal for %s timer", t.name)
			return
		case <-trigger:
			logrus.Debugf("triggered action for %s timer", t.name)
			go t.action.DoAction()
		}
	}
}

// checkTimer infinitely checks clock every `wait` duration, sends message when trigger is up
func (t *Timer) checkTimer(trigger, kill chan bool) {
	for {
		if time.Since(t.start) > t.interval {
			t.start = time.Now()
			trigger <- true
		}
		select {
		case <-kill:
			return
		default:
			time.Sleep(t.frequency)
		}
	}
}
