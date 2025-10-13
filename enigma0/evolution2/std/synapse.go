package std

import (
	"sync"
	"time"

	"git.ignitelabs.net/janos/core"
)

// Synapse represents the basic synapse type.  As synapses are rife with nuance and subtlety, to create a synapse,
// please see the std.Signal factory
type Synapse chan<- any

func NewSynapse[T any](named string, frequency float64, period time.Duration, neurons ...Neural[T]) Synapse {
	return NewSynapseRef(named, core.Ref(frequency), core.Ref(period), neurons...)
}
func NewSynapseRef[T any](named string, frequency *float64, period *time.Duration, neurons ...Neural[T]) Synapse {
	// Neurons added to a synapse are impulsed through sync.Cond
	// This happens by creating a goroutine for each neuron (which neurons can be signaled in!) that waits on the same sync.Cond as everyone else
	// Synapses are ALWAYS burst fired like this, but the neurons can intelligently activate their potentials using phase data

	// Signal.Impulse(...Thoughts)
	// Signal.Spark(...Thoughts)
	// Signal.Mute()
	// Signal.Unmute()
	// Signal.Query() // name, genesis, frequency, period, etc...

	// Signal.Tune
	// Signal.Tune.Frequency()
	// Signal.Tune.Period()
	// Signal.Tune.All(Frequency, Period)
}
