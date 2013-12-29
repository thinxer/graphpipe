package graphpipe

import (
	"time"
)

// A source that wakes up on interval.
type Timer struct {
	tid      int
	interval time.Duration
}

type TimerConfig struct {
	Interval int
}

func (i *Timer) Update(tid int) bool {
	time.Sleep(i.interval)
	i.tid = tid
	return true
}

func (i *Timer) Value() int {
	return i.tid
}

func (l *Timer) Closed() bool {
	return false
}

func NewTimer(config *TimerConfig) (*Timer, error) {
	return &Timer{interval: time.Duration(config.Interval) * time.Second}, nil
}

func init() {
	Regsiter("Timer", NewTimer)
}
