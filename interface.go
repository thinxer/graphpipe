package graphpipe

type Result int

const (
	Skip   = 0         // Emit nothing
	Update = 1 << iota // Activate children
	More   = 1 << iota // Request for next update
)

func (r Result) String() string {
	if r == 0 {
		return "-"
	}
	s := ""
	if r&Update > 0 {
		s = s + "U"
	} else if r&More > 0 {
		s = s + "+"
	}
	return s
}

// Node represents an updatable node in the pipeline.
type Node interface {
	Update(tid int) Result
}

// SourceNode can emit signals to request an update.
type SourceNode interface {
	Node
	// Start will be called on first iteration.
	// When the source is ready to emit a value, send a bool to ch.
	Start(ch chan<- bool)
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
