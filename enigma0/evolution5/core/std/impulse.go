package std

import (
	"time"

	"git.ignitelabs.net/janos/core/std"
)

type Impulse struct {
	id uint64

	// Synapse represents the originating synapse that created this Impulse.
	//Synapse

	// Heart provides a callback point to signal continued activity during long-running operations.
	//
	// For instance, if running a large batched process, you could call Heart.Beat() during every
	// cycle of your loop to give feedback of your continued execution.
	//*Heart

	// Arguments cache the impulse chain's lifetime context.
	Arguments []any

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
	Timeline *std.TemporalBuffer[Impulse]
}
