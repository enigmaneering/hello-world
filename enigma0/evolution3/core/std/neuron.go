package std

import (
	"sync"

	"git.ignitelabs.net/janos/core/std"
)

type Neuron struct {
	std.Entity
	action    func(Impulse)
	potential func(Impulse) bool
	cleanup   func(Impulse, *sync.WaitGroup)
	heart     *Heart

	created bool
}

func (n Neuron) sanityCheck() {
	if !n.created {
		panic("std.Neuron[T] must be created through std.NewNeuron[T] - otherwise, please implement the std.Neural[T] interface")
	}
}

func (n Neuron) Named() string {
	n.sanityCheck()
	return n.Name
}

func (n Neuron) Action(imp Impulse) {
	n.sanityCheck()
	if n.action != nil {
		n.action(imp)
	}
}

func (n Neuron) Potential(imp Impulse) bool {
	n.sanityCheck()
	if n.potential == nil {
		return true
	}
	return n.potential(imp)
}

func (n Neuron) Cleanup(imp Impulse, wg *sync.WaitGroup) {
	n.sanityCheck()
	if n.cleanup != nil {
		n.cleanup(imp, wg)
	} else {
		wg.Done()
	}
}

func (n Neuron) Heart() *Heart {
	return n.heart
}

func NewNeuron(named string, action func(Impulse), potential func(Impulse) bool, cleanup ...func(Impulse, *sync.WaitGroup)) Neuron {
	return Neuron{
		Entity:    std.NewEntityNamed(named),
		action:    action,
		potential: potential,
		cleanup: func(imp Impulse, wg *sync.WaitGroup) {
			for _, clean := range cleanup {
				wg.Add(1)
				clean(imp, wg)
			}
		},
		created: true,
	}
}
