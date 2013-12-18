package graphpipe_example

import (
	"fmt"
	pipe "github.com/thinxer/graphpipe"
)

type IntLogger struct {
	tid    int
	value  int
	name   string
	silent bool
	source pipe.IntSource
}

type IntLoggerConfig struct {
	Name   string
	Silent bool
}

func (f *IntLogger) Update(tid int) bool {
	f.tid, f.value = f.source.Value()
	if !f.silent {
		fmt.Printf("%s[%d]: %d[%d]\n", f.name, tid, f.value, f.tid)
	}
	return true
}

func (f *IntLogger) Value() (int, int) {
	return f.tid, f.value
}

func (f *IntLogger) Closed() bool {
	return f.source.Closed()
}

func NewIntLogger(config *IntLoggerConfig, source pipe.IntSource) *IntLogger {
	return &IntLogger{source: source, name: config.Name, silent: config.Silent}
}

func init() {
	pipe.Regsiter("IntLogger", NewIntLogger)
}
