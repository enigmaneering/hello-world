// Package relationship provides several ways of describing a psychological relationship.Constraint between a hosting std.Idea and Other.
//
// See Constraint, Open, Inclusive, and Exclusive
package relationship

// A Constraint defines several ways of psychologically restricting Other's access to a hosting std.Idea
//
// See Constraint, Open, Inclusive, and Exclusive
type Constraint byte

const (
	// Open indicates that Other can read from or write to the hosting std.Idea
	//
	// See Constraint, Open, Inclusive, and Exclusive
	Open Constraint = iota

	// Inclusive indicates that Other can read from the hosting std.Idea but can only write to it with a code
	//
	// See Constraint, Open, Inclusive, and Exclusive
	Inclusive

	// Exclusive indicates that Other can only read from or write to the hosting std.Idea with a code
	//
	// See Constraint, Open, Inclusive, and Exclusive
	Exclusive
)
