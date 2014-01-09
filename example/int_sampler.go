package graphpipe_example

import (
	"fmt"
	"log"

	pipe "github.com/thinxer/graphpipe"
)

// An int sampler
type IntSampler struct {
	tid    int
	sample int

	count    int
	interval int
	source   pipe.IntSource
}

type IntSamplerConfig struct {
	Interval int
}

func (f *IntSampler) Update(tid int) (updated bool) {
	if f.count == 0 {
		f.tid, f.sample = f.source.Value()
		updated = true
	}
	f.count = (f.count + 1) % f.interval
	return
}

func (f *IntSampler) Value() (int, int) {
	return f.tid, f.sample
}

func (f *IntSampler) Closed() bool {
	return f.source.Closed()
}

func (i *IntSampler) SetInput(source pipe.IntSource) {
	log.Println("setting input")
	i.source = source
}

func NewIntSampler(config *IntSamplerConfig) (*IntSampler, error) {
	if config.Interval <= 0 {
		return nil, fmt.Errorf("interval must be positive")
	}
	return &IntSampler{count: 0, interval: config.Interval}, nil
}

func init() {
	pipe.Regsiter("IntSampler", NewIntSampler)
}
