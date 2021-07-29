package timer_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ntbloom/rainbase/pkg/timer"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type fakeAction struct {
	counter int
	sync.Mutex
}

func (f *fakeAction) DoAction() {
	f.Lock()
	f.counter++
	f.Unlock()
}

// TestTimer basic timer should increment a counter every second for 5 seconds and then die
func TestTimer(t *testing.T) {
	fake := &fakeAction{counter: 0}
	countTimer := timer.NewTimer(time.Second, time.Millisecond*200, fake)
	var count int

	go countTimer.Loop()
	time.Sleep(time.Millisecond * 200)

	fake.Lock()
	count = fake.counter
	fake.Unlock()
	assert.Equal(t, 0, count)

	time.Sleep(time.Second * 5)
	countTimer.Kill <- true

	fake.Lock()
	count = fake.counter
	fake.Unlock()

	logrus.Infof("count=%d", count)
	assert.Equal(t, count, 5)
}
