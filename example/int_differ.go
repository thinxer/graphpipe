package graphpipe_example

import (
	pipe "github.com/thinxer/graphpipe"
)

type IntDiffer struct {
	tid   int
	value int

	a pipe.IntSource
	b pipe.IntSource
}

func (f *IntDiffer) Update(tid int) bool {
	if f.Closed() {
		return false
	}
	_, val1 := f.a.Value()
	_, val2 := f.b.Value()
	f.tid, f.value = tid, val1-val2
	return true
}

func (f *IntDiffer) Value() (int, int) {
	return f.tid, f.value
}

func (f *IntDiffer) Closed() bool {
	return f.a.Closed() || f.b.Closed()
}

func NewIntDiffer(config *struct{}, a, b pipe.IntSource) (*IntDiffer, error) {
	return &IntDiffer{a: a, b: b}, nil
}

func init() {
	pipe.Regsiter("IntDiffer", NewIntDiffer)
}
