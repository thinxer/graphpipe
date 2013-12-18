package graphpipe_example

import (
	pipe "github.com/thinxer/graphpipe"
)

// This delayer is only for demo.
// The 'cache' field will leak memory. You may consider a circular queue instead.
type IntDelayer struct {
	tid   int
	value int

	cache []struct {
		tid, value int
	}
	count  int
	delay  int
	source pipe.IntSource
}

type IntDelayerConfig struct {
	Delay int
}

func (f *IntDelayer) Update(_ int) (updated bool) {
	if f.count == f.delay {
		if len(f.cache) > 0 {
			updated, f.tid, f.value = true, f.cache[0].tid, f.cache[0].value
			f.cache = f.cache[1:len(f.cache)]
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
		return len(f.cache) == 0
	}
}

func NewIntDelayer(config *IntDelayerConfig, source pipe.IntSource) *IntDelayer {
	if config.Delay <= 0 {
		panic("delay must be positive: you cannot travel to the future!")
	}
	return &IntDelayer{delay: config.Delay, source: source}
}

func init() {
	pipe.Regsiter("IntDelayer", NewIntDelayer)
}
