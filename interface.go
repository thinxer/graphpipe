package graphpipe

// An updatable node in the pipeline
type Node interface {
	// Return true to activate nodes depending on this one.
	Update(tid int) (updated bool)
	// Return true if the node won't output anything anymore
	// This method is usually also required by the Source interfaces.
	Closed() bool
}

// A source emitting integers
type IntSource interface {
	Value() (int, int)
	Closed() bool
}

// A source emitting integers
type Int32Source interface {
	Value() (int, int32)
	Closed() bool
}

// A source emitting integers
type Int64Source interface {
	Value() (int, int64)
	Closed() bool
}

// A source emitting strings
type StringSource interface {
	Value() (int, string)
	Closed() bool
}

// A source emitting float32s
type Float32Source interface {
	Value() (int, float32)
	Closed() bool
}

// A source emitting float64s
type Float64Source interface {
	Value() (int, float64)
	Closed() bool
}

// A source emitting nothing
type NilSource interface {
	Value() int
	Closed() bool
}

// A source emitting anything
type AnySource interface {
	Closed() bool
}
