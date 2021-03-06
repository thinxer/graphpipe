package graphpipe_example

import (
	"fmt"

	pipe "github.com/thinxer/graphpipe"
)

type IntDelayer struct {
	tid   int
	value int

	cache  []int
	delay  int
	closed bool
	source pipe.IntSource
}

type IntDelayerConfig struct {
	Delay int
}

func (f *IntDelayer) Update(tid int) (result pipe.Result) {
	if f.delay == 0 {
		if len(f.cache) > 0 {
			f.tid, f.value = tid, f.cache[0]
			f.cache = f.cache[1:]
			result |= pipe.Update
			if f.source.Closed() {
				result |= pipe.More
			}
		}
	} else {
		f.delay--
	}

	stid, val := f.source.Value()
	if stid == tid {
		f.cache = append(f.cache, val)
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
	pipe.Register("IntDelayer", NewIntDelayer)
}
