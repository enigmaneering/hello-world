package std

import (
	"fmt"

	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution2/core/enum/relationship"
)

// A Disclosure provides a control surface over access to an Idea's Thought.  Whoever holds a reference
// to the disclosure can control psychological access to the idea - either by applying a Stringable code
// or by directly setting the relationship.Constraint.
//
// NOTE: The disclosure code is an 'any' type to allow lazy comparison of a live structure at runtime.
type Disclosure struct {
	// Constraint defines the psychological relationship applied to this Idea
	relationship.Constraint
	code any
}

func (d *Disclosure) sanityCheck() {
	if !Stringable(d.code) {
		panic(fmt.Errorf("the std.Disclosure's code is a %T - which is not a stringable type", d.code))
	}
}

// Code sets and/or gets the Disclosure code.  If no value is provided, this just returns the code - if values
// are provided, the first value is used as the disclosure code before returning it.
//
// NOTE: This will panic if not provided a Stringable type!
func (d *Disclosure) Code(value ...any) any {
	if len(value) > 0 && Stringable(value[0]) {
		d.code = value[0]
	}
	return value
}

// Check validates if the provided code matches the Disclosure code.
//
// NOTE: This will panic if not provided a Stringable type!
func (d *Disclosure) Check(code any) bool {
	return Stringify(code) == Stringify(d.code)
}
