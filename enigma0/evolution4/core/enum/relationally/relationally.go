// Package relationally provides several ways of describing a psychological relationally.Constrained between a hosting std.Idea and Other.
//
// See Constrained, Open, Inclusive, and Exclusive
package relationally

// Constrained defines several ways of filtering Other's access to a std.Thought
//
// See Constrained, Open, Inclusive, and Exclusive
type Constrained byte

const (
	// Open indicates that Other can indiscriminately read from or write to the std.Thought
	//
	// See Constrained, Open, Inclusive, and Exclusive
	Open Constrained = iota

	// Inclusive indicates that Other can read from the std.Thought but can only write to it with a code
	//
	// See Constrained, Open, Inclusive, and Exclusive
	Inclusive

	// Exclusive indicates that Other can only read from or write to the std.Thought with a code
	//
	// See Constrained, Open, Inclusive, and Exclusive
	Exclusive
)
