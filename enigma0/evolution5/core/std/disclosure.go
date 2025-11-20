package std

import (
	"fmt"

	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution5/core/enum/relationally"
)

// A Disclosure describes the codified conditions and constraints behind Others' access to a Thought.
type Disclosure struct {
	// Constraint defines the psychological relationship applied to this Thought
	Constraint relationally.Constrained
	code       any
}

func (d *Disclosure) sanityCheck(codes ...any) {
	if !negotiable(d.code) {
		if !Stringable(d.code) {
			panic(fmt.Errorf("the std.Disclosure's code is a %T - which is not a std.Stringable or std.Negotiable type", d.code))
		}
	}
	if !negotiable(codes...) {
		if !Stringable(codes...) {
			panic(fmt.Errorf("the provided code is a %T - which is not a std.Stringable or std.Negotiable type", codes[0]))
		}
	}
}

// Code sets and/or gets the Disclosure code.  If no value is provided, this just returns the code - if values
// are provided, the first value is set as the disclosure code before returning it.
//
// NOTE: This will panic if provided anything but a Stringable or Negotiable type!
func (d *Disclosure) Code(value ...any) any {
	if len(value) > 0 {
		if negotiable(value[0]) || Stringable(value[0]) {
			d.code = value[0]
		} else {
			panic(fmt.Errorf("the provided code is a %T - which is not a std.Stringable or std.Negotiable type", value[0]))
		}
	}
	return value
}

// Check validates the provided code.
//
// 0. If a nil (or absent) code is provided, a string comparison of "" is made
//
// 1. If a Negotiable code is provided, the disclosure.code.Negotiate(code) must yield "true".
//
// 2. All other provided codes must be Stringable, and a string comparison is performed for equivalency
//
// 3. Otherwise, this will panic during its "sanity check"
func (d *Disclosure) Check(code ...any) bool {
	d.sanityCheck()

	if len(code) == 0 && (d.code == nil || d.code == "") {
		return true
	} else if len(code) > 0 {
		return Stringify(code[0]) == Stringify(d.code)
	}
	return false
}
