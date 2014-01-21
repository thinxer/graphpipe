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

func (f *IntDiffer) Update(tid int) pipe.UpdateResult {
	_, val1 := f.a.Value()
	_, val2 := f.b.Value()
	f.tid, f.value = tid, val1-val2
	return pipe.Updated
}

func (f *IntDiffer) Value() (int, int) {
	return f.tid, f.value
}

func (f *IntDiffer) Closed() bool {
	return f.a.Closed() && f.b.Closed()
}

func (i *IntDiffer) SetInput(a, b pipe.IntSource) {
	i.a, i.b = a, b
}

func NewIntDiffer(config *struct{}) (*IntDiffer, error) {
	return &IntDiffer{}, nil
}

func init() {
	pipe.Regsiter("IntDiffer", NewIntDiffer)
}
