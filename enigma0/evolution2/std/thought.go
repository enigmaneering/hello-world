package std

import (
	"strconv"
	"sync"

	"git.ignitelabs.net/janos/core/sys/id"
)

type Thought[T any] struct {
	Path
	revelation T
	gate       *sync.Mutex
	created    bool
}

// NewThought creates a new instance of a Thought[T] and assigns it a unique Path identifier by calling id.Next()
//
// NOTE: If you'd like to assign the thought's path directly, you may provide it through the variadic.
func NewThought[T any](revelation T, path ...string) Thought[T] {
	p := []string{strconv.FormatUint(id.Next(), 10)}
	if len(path) > 0 {
		p = path
	}
	return Thought[T]{
		Path:       p,
		revelation: revelation,
		gate:       new(sync.Mutex),
		created:    true,
	}
}

func (t *Thought[T]) sanityCheck() {
	if !t.created {
		panic("std.Thought[T] must be created through std.NewThought[T]")
	}
	if t.gate == nil {
		t.gate = new(sync.Mutex)
	}
}

// Revelation sets and/or gets the Thought's inner revelation.  If no value is
// provided, it merely gets - otherwise it sets the value before returning it.
func (t *Thought[T]) Revelation(value ...T) T {
	t.sanityCheck()
	t.gate.Lock()
	defer t.gate.Unlock()

	if len(value) > 0 {
		t.revelation = value[0]
	}
	return t.revelation
}
