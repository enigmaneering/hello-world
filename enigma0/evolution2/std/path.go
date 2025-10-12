package std

import "strings"

// A Path is a sequence of string identifiers which can act as a chain of custody.
type Path []string

// Emit outputs the path as a delimited string.
//
// NOTE: If no delimiter is provided, ⇝ is used - otherwise, only the first delimiter provided is used.
func (t Path) Emit(delimiter ...string) string {
	d := "⇝"
	if len(delimiter) > 0 {
		d = delimiter[0]
	}

	return strings.Join(t, d)
}
