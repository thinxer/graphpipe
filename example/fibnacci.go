package graphpipe_example

import (
	p "github.com/thinxer/graphpipe"
)

// A sample data source

type Fibonacci struct {
	tid   int
	a, b  int
	limit int
}

type FibonacciConfig struct {
	Seed1, Seed2 int
	Limit        int
}

func newFibonacci(config *FibonacciConfig) (*Fibonacci, error) {
	return &Fibonacci{a: config.Seed1, b: config.Seed2, limit: config.Limit}, nil
}

func (f *Fibonacci) Start(ch chan bool) {
	ch <- true
	close(ch)
}

func (f *Fibonacci) Update(tid int) p.Result {
	if f.limit > 0 {
		f.a, f.b, f.tid = f.b, f.a+f.b, tid
		f.limit--
		return p.Update | p.More
	} else {
		f.limit--
		return p.Update
	}
}

func (f *Fibonacci) Value() (int, int) {
	return f.tid, f.a
}

func (f *Fibonacci) Closed() bool {
	return f.limit < 0
}

func init() {
	p.Register("Fibonacci", newFibonacci)
}
