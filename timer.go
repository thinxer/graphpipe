package graphpipe

import "time"

// A source that wakes up on interval.
type Timer struct {
	tid      int
	interval time.Duration
	source   AnySource
	closing  chan bool
	closed   bool
}

type TimerConfig struct {
	Interval int
}

func (i *Timer) Start(ch chan bool) {
L:
	for {
		select {
		case <-i.closing:
			break L
		case <-time.After(i.interval):
			ch <- true
		}
	}
	close(ch)
}

func (i *Timer) Update(tid int) UpdateResult {
	if i.closed {
		return Skip
	}
	if i.source.Closed() {
		i.closed = true
		i.closing <- true
		return Skip
	}
	i.tid = tid
	return Updated
}

func (i *Timer) Value() int {
	return i.tid
}

func (i *Timer) Closed() bool {
	return i.closed
}

func newTimer(config *TimerConfig, source AnySource) (*Timer, error) {
	return &Timer{interval: time.Duration(config.Interval) * time.Second, closing: make(chan bool), source: source}, nil
}

func init() {
	Regsiter("Timer", newTimer)
}
