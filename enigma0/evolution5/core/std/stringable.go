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
// - any of the tiny.Numeric constraint members, including:
//
// - - big.Int - using Text(10)
//
// - - big.Float - using Text("f", atlas.Precision)
//
// - - big.Rat - using String()
//
// - - float32/float64 - using strconv.FormatFloat(fmt:"f", prec:-1)
//
// See Stringable, Stringify, and StringifyMany
func Stringable(values ...any) bool {
	allGood := true
	for _, value := range values {
		switch value.(type) {
		case nil, string, fmt.Stringer:
		default:
			_, err := num.ToStringSafe(value)
			if err != nil {
				allGood = false
			}
		}
	}
	return allGood
}
