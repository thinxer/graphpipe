package graphpipe_example

import (
	"fmt"

	pipe "github.com/thinxer/graphpipe"
)

type IntLogger struct {
	name   string
	silent bool
	source pipe.IntSource
}

type IntLoggerConfig struct {
	Name   string
	Silent bool
}

func (f *IntLogger) Update(tid int) pipe.UpdateResult {
	stid, value := f.source.Value()
	if !f.silent {
		fmt.Printf("%s[%d]: %d[%d]\n", f.name, tid, value, stid)
	}
	return pipe.Updated
}

func (f *IntLogger) Closed() bool {
	return f.source.Closed()
}

func (i *IntLogger) SetInput(source pipe.IntSource) {
	i.source = source
}

func NewIntLogger(config *IntLoggerConfig) (*IntLogger, error) {
	return &IntLogger{name: config.Name, silent: config.Silent}, nil
}

func init() {
	pipe.Regsiter("IntLogger", NewIntLogger)
}
