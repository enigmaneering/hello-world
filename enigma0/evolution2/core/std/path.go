package std

import (
	"strings"
)

// A Path is a sequence of Stringable identifiers which can be used to locate something.
type Path []any

func (p Path) sanityCheck() {
	StringifyMany(p...)
}

// String outputs the path's components as a '⇝' delimited string.
func (p Path) String() string {
	return p.Emit()
}

// Emit outputs the path as a delimited string.
//
// NOTE: If no delimiter is provided, '⇝' is used - otherwise, only the first delimiter provided is used.
func (p Path) Emit(delimiter ...string) string {
	d := "⇝"
	if len(delimiter) > 0 {
		d = delimiter[0]
	}

	return strings.Join(StringifyMany(p), d)
}
