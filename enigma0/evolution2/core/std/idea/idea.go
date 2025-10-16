package idea

import (
	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution2/core/enum/relationship"
	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution2/core/std"
)

var origin std.Nexus

// New places a thought at the provided path relative to the origin nexus and returns a *std.Disclosure controlling its access.
// The path could be a std.Stringable or a std.Pathable, such as a std.Idea.
//
// NOTE: Thoughts can still be shared amongst ideas with different disclosed access constraints! =)
//
// See New, Reveal, and Describe
func New[T any](revelation T, path any, disclosure ...*std.Disclosure) *std.Disclosure {
	var p std.Path
	switch typed := path.(type) {
	case std.Path:
		p = typed
	case std.Pathable:
		p = typed.Path()
	default:
		if std.Stringable(path) {
			p = std.Path{std.Stringify(path)}
		} else {
			panic("the provided path is not a std.Stringable or std.Pathable type")
		}
	}

	d := &std.Disclosure{
		Constraint: relationship.Open,
		Code:       "",
	}
	if len(disclosure) > 0 {
		d = disclosure[0]
	}

	id := std.NewIdea[T](std.NewThought(revelation), p, d)
	// TODO: Place on the nexus
	return d
}

// Reveal follows the provided path towards a revelation.  The path could be a std.Stringable, std.Path, or std.Idea.
//
// See New, Reveal, and Describe
func Reveal[T any](path any, code ...any) (T, error) {
	var zero T
	return zero, nil
}

// Describe sets the revelation at the end of the provided path.  The path could be a std.Stringable, std.Path, or std.Idea.
//
// See New, Reveal, and Describe
func Describe[T any](revelation T, path any, code ...any) error {
	return nil
}
