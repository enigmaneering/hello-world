// Package relationally provides several ways of describing a psychological relationally.Constrained between a hosting std.Idea and Other.
//
// See Constrained, Open, Inclusive, and Exclusive
package relationally

// A Constrained defines several ways of psychologically restricting Other's access to a hosting std.Idea
//
// See Constrained, Open, Inclusive, and Exclusive
type Constrained byte

const (
	// Open indicates that Other can read from or write to the hosting std.Idea
	//
	// See Constrained, Open, Inclusive, and Exclusive
	Open Constrained = iota

	// Inclusive indicates that Other can read from the hosting std.Idea but can only write to it with a code
	//
	// See Constrained, Open, Inclusive, and Exclusive
	Inclusive

	// Exclusive indicates that Other can only read from or write to the hosting std.Idea with a code
	//
	// See Constrained, Open, Inclusive, and Exclusive
	Exclusive
)
