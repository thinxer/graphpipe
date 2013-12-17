package graphpipe

import (
	"fmt"
)

// An int sampler

type IntLogger struct {
	tid    int
	value  int
	name   string
	silent bool
	source IntSource
}

type IntLoggerConfig struct {
	Name   string
	Silent bool
}

func (f *IntLogger) Update(tid int) bool {
	f.tid, f.value = f.source.Value()
	if !f.silent {
		fmt.Printf("%s[%d]: %d\n", f.name, f.tid, f.value)
	}
	return true
}

func (f *IntLogger) Value() (int, int) {
	return f.tid, f.value
}

func (f *IntLogger) Closed() bool {
	return f.source.Closed()
}

func NewIntLogger(config *IntLoggerConfig, source IntSource) *IntLogger {
	return &IntLogger{source: source, name: config.Name, silent: config.Silent}
}

func init() {
	Regsiter("IntLogger", NewIntLogger)
}
