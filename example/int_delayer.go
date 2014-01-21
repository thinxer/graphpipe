package graphpipe_example

import (
	"fmt"

	pipe "github.com/thinxer/graphpipe"
)

type IntDelayer struct {
	tid   int
	value int

	cache []struct {
		tid, value int
	}
	count  int
	delay  int
	closed bool
	source pipe.IntSource
}

type IntDelayerConfig struct {
	Delay int
}

func (f *IntDelayer) Update(_ int) (updated pipe.UpdateResult) {
	if f.count == f.delay {
		if len(f.cache) > 0 {
			updated, f.tid, f.value = pipe.Updated, f.cache[0].tid, f.cache[0].value
			f.cache = f.cache[1:]
			if f.source.Closed() {
				updated = pipe.HasMore
			}
		} else {
			f.closed = true
		}
	} else {
		f.count++
	}

	if !f.source.Closed() {
		tid, val := f.source.Value()
		f.cache = append(f.cache, struct{ tid, value int }{tid, val})
	}

	return
}

func (f *IntDelayer) Value() (int, int) {
	return f.tid, f.value
}

func (f *IntDelayer) Closed() bool {
	if f.cache == nil {
		return f.source.Closed()
	} else {
		return f.closed
	}
}

func (i *IntDelayer) SetInput(source pipe.IntSource) {
	i.source = source
}

func NewIntDelayer(config *IntDelayerConfig) (*IntDelayer, error) {
	if config.Delay <= 0 {
		return nil, fmt.Errorf("delay must be positive: you cannot travel to the future!")
	}
	return &IntDelayer{delay: config.Delay}, nil
}

func init() {
	pipe.Regsiter("IntDelayer", NewIntDelayer)
}
