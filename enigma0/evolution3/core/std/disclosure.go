package std

import (
	"fmt"

	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution3/core/enum/relationally"
)

// A Disclosure provides a control surface over access to an Idea's Thought.  Whoever holds a reference
// to the disclosure can control psychological access to the idea - either by applying a Stringable code
// or by making it relationally.Constrained.
//
// NOTE: The disclosure code is an 'any' type to allow lazy comparison of a live structure at runtime.
type Disclosure struct {
	// Constraint defines the psychological relationship applied to this Idea
	Constraint relationally.Constrained
	code       any
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
// NOTE: nil or an absent code is equivalent to an empty string and only the first code value is considered.
//
// NOTE: This will panic if not provided a Stringable type!
func (d *Disclosure) Check(code ...any) bool {
	if len(code) == 0 && (d.code == nil || d.code == "") {
		return true
	} else if len(code) > 0 {
		return Stringify(code[0]) == Stringify(d.code)
	}
	return false
}
