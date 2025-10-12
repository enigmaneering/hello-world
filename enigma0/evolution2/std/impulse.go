package std

import (
	"time"

	"git.ignitelabs.net/janos/core/std"
)

type Impulse[T any] struct {
	id uint64

	// Synapse represents the originating synapse that created this Impulse.
	Synapse

	// Thought represents the underlying revelation provided with this Impulse.
	Thought[T]

	// Bridge holds the chain of named activations in creating this Impulse.
	Bridge Path

	// Epoch holds the moment of impulse coordination, typically the synapse's genesis moment.
	Epoch time.Time

	// Inception holds the moment this impulse was created.
	Inception time.Time

	// Activated holds the moment this impulse's action was fired.
	//
	// NOTE: This will remain nil until the action's potential returns "high."
	Activated *time.Time

	// Completed holds the moment this impulse's action completed.
	//
	// NOTE: This will remain nil until the action finishes its current execution.
	Completed *time.Time

	// Timeline holds a temporal buffer of prior synaptic activations.
	//
	// NOTE: The impulse is not added to the buffer before calling the potential or action.
	Timeline *std.TemporalBuffer[Impulse[T]]
}
