package graphpipe_example

import (
	"github.com/thinxer/graphpipe"
)

// A sample data source

type Fibonacci struct {
	a, b  int
	tid   int
	count int
	limit int
}

type FibonacciConfig struct {
	Seed1, Seed2 int
	Limit        int
}

func (f *Fibonacci) Update(tid int) bool {
	if f.count < f.limit {
		f.a, f.b = f.b, f.a+f.b
		f.tid = tid
		f.count++
		return true
	} else if f.count == f.limit {
		f.count = -1
		return false
	}
	return false
}

func (f *Fibonacci) Closed() bool {
	return f.count < 0
}

func (f *Fibonacci) Value() (int, int) {
	return f.tid, f.a
}

func newFibonacci(config *FibonacciConfig) (*Fibonacci, error) {
	return &Fibonacci{a: config.Seed1, b: config.Seed2, count: 0, limit: config.Limit}, nil
}

func init() {
	graphpipe.Regsiter("Fibonacci", newFibonacci)
}
