package std

import (
	"fmt"

	"git.ignitelabs.net/janos/core/sys/num"
)

// Stringify converts the provided value into a string.  If the value does not satisfy Stringable, this will panic.
//
// NOTE: For a runtime-safe check, see Stringable
//
// See Stringable, StringableMany, Stringify, and StringifyMany
func Stringify(value any) string {
	switch raw := value.(type) {
	case string:
		return raw
	case fmt.Stringer:
		return raw.String()
	default:
		out, err := num.ToStringSafe(value)
		if err != nil {
			panic(fmt.Errorf("%T is not a stringable type", raw))
		}
		return out
	}
}

// StringifyMany converts the provided values into a []string.  If the value does not satisfy Stringable, this will panic.
//
// NOTE: For a runtime-safe check, see StringableMany
//
// See Stringable, StringableMany, Stringify, and StringifyMany
func StringifyMany(values ...any) []string {
	out := make([]string, len(values))
	for i, raw := range values {
		out[i] = Stringify(raw)
	}
	return out
}
