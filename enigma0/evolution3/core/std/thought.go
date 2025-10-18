package std

import (
	"sync"
)

// A Thought is a gated runtime object.  It holds a value of type T accessible through a Revelation method,
// which performs all read and write operations through a mutex.
//
// NOTE: This is a value type but can hold a reference =)
type Thought[T any] struct {
	revelation T
	gate       *sync.Mutex
	created    bool
}

// NewThought creates a new instance of a Thought[T]
func NewThought[T any](revelation T) Thought[T] {
	return Thought[T]{
		revelation: revelation,
		gate:       new(sync.Mutex),
		created:    true,
	}
}

func (t Thought[T]) sanityCheck() {
	if !t.created {
		panic("a std.Thought[T] must be created through std.NewThought[T]")
	}
	if t.gate == nil {
		t.gate = new(sync.Mutex)
	}
}

// Revelation sets and/or gets the Thought's inner revelation.  If no value is
// provided, it merely gets - otherwise it sets the value before returning it.
func (t Thought[T]) Revelation(value ...T) T {
	t.sanityCheck()
	t.gate.Lock()
	defer t.gate.Unlock()

	if len(value) > 0 {
		t.revelation = value[0]
	}
	return t.revelation
}
