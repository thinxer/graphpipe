package graphpipe

type UpdateResult int

const (
	Skip UpdateResult = iota
	Updated
	HasMore
	End
)

// Node represents an updatable node in the pipeline.
type Node interface {
	// Return true to activate nodes depending on this one.
	Update(tid int) (updated UpdateResult)
	// Return true if the node won't output anything anymore
	// This method is usually also required by the Source interfaces.
	Closed() bool
}

// SourceNode can emit values.
type SourceNode interface {
	Node
	// Start will be called on first iteration.
	// When the source is ready to emit a value, send a bool to ch.
	Start(ch chan bool)
}

// IntSouce emits int.
type IntSource interface {
	Value() (int, int)
	Closed() bool
}

// Int32Source emits int32.
type Int32Source interface {
	Value() (int, int32)
	Closed() bool
}

// Int64Source emits int64.
type Int64Source interface {
	Value() (int, int64)
	Closed() bool
}

// StringSource emits string.
type StringSource interface {
	Value() (int, string)
	Closed() bool
}

// Float32Source emits float32.
type Float32Source interface {
	Value() (int, float32)
	Closed() bool
}

// Float64Source emits float64.
type Float64Source interface {
	Value() (int, float64)
	Closed() bool
}

// NilSource emits nothing.
type NilSource interface {
	Value() int
	Closed() bool
}

// AnySource emits anything (or do no emit at all).
type AnySource interface {
	Closed() bool
}
