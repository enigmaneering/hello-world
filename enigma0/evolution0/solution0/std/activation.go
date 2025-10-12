package std

import "time"

type Activation struct {
	id uint64

	// Epoch represents the moment the synaptic connection was created.
	Epoch time.Time

	// Inception represents the moment a neuron received a synaptic event.
	Inception time.Time

	// Fired represents the moment a neuron passed it's potential.
	Fired *time.Time

	// Completion represents the moment a neuron finished execution.
	Completion *time.Time
}

func (a *Activation) ID() uint64 {
	return a.id
}
