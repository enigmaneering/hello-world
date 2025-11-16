package std

import (
	"strings"
)

// A Path is a sequence of Stringable Step points which can be used to relatively Locate something.
type Path []any

func (p Path) sanityCheck() {
	StringifyMany(p)
}

// String outputs the Path's steps as a '⇝' delimited string, minus any code information.
func (p Path) String() string {
	return strings.Join(StringifyMany(p), "⇝")
}
