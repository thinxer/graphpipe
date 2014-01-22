package graphpipe_example

import (
	p "github.com/thinxer/graphpipe"
)

// A sample data source

type Fibonacci struct {
	a, b  int
	limit int

	tid     int
	value   int
	pending chan int
	closed  bool
}

type FibonacciConfig struct {
	Seed1, Seed2 int
	Limit        int
}

func newFibonacci(config *FibonacciConfig) (*Fibonacci, error) {
	return &Fibonacci{a: config.Seed1, b: config.Seed2, limit: config.Limit, pending: make(chan int, 128)}, nil
}

func (f *Fibonacci) Start(ch chan bool) {
	for f.limit > 0 {
		f.limit--
		f.a, f.b = f.b, f.a+f.b
		f.pending <- f.a
		ch <- true
	}
	close(ch)
	close(f.pending)
}

func (f *Fibonacci) Update(tid int) p.UpdateResult {
	v, ok := <-f.pending
	if ok {
		f.tid, f.value = tid, v
		return p.Updated
	} else {
		f.closed = true
		return p.Skip
	}
}

func (f *Fibonacci) Closed() bool {
	return f.closed
}

func (f *Fibonacci) Value() (int, int) {
	return f.tid, f.value
}

func init() {
	p.Regsiter("Fibonacci", newFibonacci)
}
