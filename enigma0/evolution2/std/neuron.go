package std

import (
	"sync"

	"git.ignitelabs.net/janos/core/std"
)

type Neuron[T any] struct {
	std.Entity
	action    func(Impulse[T])
	potential func(Impulse[T]) bool
	cleanup   func(Impulse[T], *sync.WaitGroup)

	created bool
}

func (n Neuron[T]) sanityCheck() {
	if !n.created {
		panic("std.Neuron[T] must be created through std.NewNeuron[T] - otherwise, please implement the std.Neural[T] interface")
	}
}

func (n Neuron[T]) Named() string {
	n.sanityCheck()
	return n.Name
}

func (n Neuron[T]) Action(imp Impulse[T]) {
	n.sanityCheck()
	if n.action != nil {
		n.action(imp)
	}
}

func (n Neuron[T]) Potential(imp Impulse[T]) bool {
	n.sanityCheck()
	if n.potential == nil {
		return true
	}
	return n.potential(imp)
}

func (n Neuron[T]) Cleanup(imp Impulse[T], wg *sync.WaitGroup) {
	n.sanityCheck()
	if n.cleanup != nil {
		n.cleanup(imp, wg)
	} else {
		wg.Done()
	}
}

func NewNeuron[T any](named string, action func(Impulse[T]), potential func(Impulse[T]) bool, cleanup ...func(Impulse[T], *sync.WaitGroup)) Neuron[T] {
	return Neuron[T]{
		Entity:    std.NewEntityNamed(named),
		action:    action,
		potential: potential,
		cleanup: func(imp Impulse[T], wg *sync.WaitGroup) {
			for _, clean := range cleanup {
				wg.Add(1)
				clean(imp, wg)
			}
		},
		created: true,
	}
}
