package std

// Neural represents the standard methods for driving an action potential off of a signal.  Please navigate into the documentation of each member for details.
type Neural interface {
	// Named returns the underlying Entity.Given.
	Named() string

	// Signal provides the signaling channel through which driving neural activity is driven.
	Signal() chan<- any

	// Action is the neural endpoint to fire with temporal context through an Impulse reference.
	Action(*Impulse)

	// Potential takes in temporal context through an Impulse reference and should return whether or not to fire the Action.
	Potential(*Impulse) bool

	// Cleanup is guaranteed to be called once the synaptic life completes, regardless of any activation.
	Cleanup(*Impulse)
}
