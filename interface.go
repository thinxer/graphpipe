package graphpipe

// An updatable node in the pipeline
type Node interface {
	// Return true to activate its dependecies
	Update(tid int) (updated bool)
	// Return true if the node won't output anything anymore
	// This method is usually also required by the Source interfaces.
	Closed() bool
}

// Common sources
// You can define your own, of course.

// A source emitting integers
type IntSource interface {
	Value() (int, int)
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
type AnythingSource interface {
	Value() (int, interface{})
	Closed() bool
}
