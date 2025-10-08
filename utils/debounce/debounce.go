package debounce

import (
	"sync/atomic"
	"time"
)

type Debounce struct {
	state  int32
	signal chan struct{}
	delay  time.Duration
}

func NewDebounce(delay time.Duration) *Debounce {
	d := &Debounce{
		state:  0,
		signal: make(chan struct{}),
		delay:  delay,
	}
	go d.loopUnlock()
	return d
}

func (d *Debounce) TryLock() bool {
	if atomic.CompareAndSwapInt32(&d.state, 0, 1) {
		d.signal <- struct{}{}
		return true
	}
	return false
}

func (d *Debounce) IsLock() bool {
	return atomic.LoadInt32(&d.state) == 1
}

func (d *Debounce) loopUnlock() {
	for range d.signal {
		time.Sleep(d.delay)
		atomic.StoreInt32(&d.state, 0)
	}
}
