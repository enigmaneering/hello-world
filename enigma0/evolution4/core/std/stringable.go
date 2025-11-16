package std

import (
	"fmt"

	"git.ignitelabs.net/janos/core/sys/num"
)

// A Stringable value is any type that can be intuitively parsed into a string.  Currently, this accepts the following:
//
// - string
//
// - nil - interpreted as an empty string
//
// - any type that satisfies fmt.Stringer
//
// - any num.Primitive
//
// - any big.Int - using Text(10)
//
// - any big.Float - using Text("f", atlas.Precision)
//
// - any big.Rat - using String()
//
// See Stringable, StringableMany, Stringify, and StringifyMany
func Stringable(value any) bool {
	switch value.(type) {
	case nil, string, fmt.Stringer:
		return true
	default:
		_, err := num.ToStringSafe(value)
		if err != nil {
			return false
		}
		return true
	}
}

// StringableMany is a convenience method for checking if many values are Stringable - please see its documentation.
//
// See Stringable, StringableMany, Stringify, and StringifyMany
func StringableMany(values ...any) bool {
	for _, value := range values {
		if !Stringable(value) {
			return false
		}
	}
	return true
}
