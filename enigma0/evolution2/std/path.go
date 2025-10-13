package std

import (
	"fmt"
	"strings"
)

// A Path is a sequence of string identifiers which can act as a chain of custody.
// NOTE: This is a Stringer type so we will be able to perform type switching at a later time
type Path []fmt.Stringer

// Emit outputs the path as a delimited string.
//
// NOTE: If no delimiter is provided, ⇝ is used - otherwise, only the first delimiter provided is used.
func (t Path) Emit(delimiter ...string) string {
	d := "⇝"
	if len(delimiter) > 0 {
		d = delimiter[0]
	}

	strs := make([]string, len(t))
	for i, s := range t {
		strs[i] = s.String()
	}
	return strings.Join(strs, d)
}
