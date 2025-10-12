package std

import (
	"sync"
	"time"

	"git.ignitelabs.net/janos/core/std"
)

type Neural[T any] interface {
	Named() string

	Action(Impulse[T])

	Potential(Impulse[T]) bool

	Cleanup(Impulse[T])
}

// Synapse represents the basic synapse type.  As synapses are rife with nuance and subtlety, to create a synapse,
// please see the std.Signal factory
type Synapse chan<- any

type Thought[T any] struct {
	Realization T
	gate        *sync.Mutex
}

func (t Thought[T]) sanityCheck() {
	if t.gate == nil {
		t.gate = new(sync.Mutex)
	}
}

func (t Thought[T]) Lock() {
	t.sanityCheck()
	t.gate.Lock()
}

func (t Thought[T]) Unlock() {
	t.sanityCheck()
	t.gate.Unlock()
}

type Impulse[T any] struct {
	std.Entity
	Synapse
	Thought[T]

	Bridge   []string
	Creation time.Time
	Epoch    time.Time
	Fire     time.Time

	Completed *time.Time

	Timeline *std.TemporalBuffer[Thought[T]]
}

type signalMaker[T any] byte

func Signal[T any]() signalMaker[T] {
	return signalMaker[T](0)
}

type impulse[T any] []*Thought[T]

func (signalMaker[T]) Impulse(thoughts ...*Thought[T]) any {
	return thoughts
}

type spark[T any] []*Thought[T]

func (signalMaker[T]) Spark(thoughts ...*Thought[T]) any {
	return thoughts
}
