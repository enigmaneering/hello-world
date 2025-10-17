package std

import (
	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution2/core/enum/errs"
	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution2/core/enum/relationally"
)

// An Idea represents a locatable and psychologically constrainable Thought.
//
// NOTE: Please use the idea.New, idea.Reveal, idea.Describe, and idea.Decay functions to orchestrate the std.Nexus of ideas.
type Idea[T any] struct {
	path       Path
	thought    *Thought[T]
	disclosure *Disclosure
}

func NewIdea[T any](thought *Thought[T], path Path, disclosure *Disclosure) Idea[T] {
	return Idea[T]{
		path:       path,
		thought:    thought,
		disclosure: disclosure,
	}
}

func (id *Idea[T]) sanityCheck() {
	if id.path == nil {
		id.path = Path{""}
	}
	id.path.sanityCheck()
	id.thought.sanityCheck()
}

func (id *Idea[T]) Path() Path {
	return id.path
}

// Reveal returns the underlying Thought's revelation.
// If the Idea is relationally.Exclusive and an invalid code is provided, an errs.CodeInvalid will be returned.
// Otherwise, codeless revelation is always permitted for relationally.Inclusive and relationally.Open constraints.
func (id *Idea[T]) Reveal(code ...any) (T, error) {
	if id.disclosure.Constraint != relationally.Exclusive || id.disclosure.Check(code...) {
		return id.thought.Revelation(), nil
	}
	var zero T
	return zero, errs.CodeInvalid
}

// Describe sets the underlying Thought's revelation.
// If the Idea is not relationally.Open and the provided code is invalid, an errs.CodeInvalid will be returned.
func (id *Idea[T]) Describe(revelation T, code ...any) error {
	if id.disclosure.Constraint == relationally.Open || id.disclosure.Check(code...) {
		id.thought.Revelation(revelation)
	}
	return errs.CodeInvalid
}

// Relative navigates the Idea using the provided Path and the following rules:
//
// 0 - If the path cannot be walked, an errs.InvalidPath structure is returned.
//
// 1 - If the path is empty, the Idea's Thought is returned.
//
// 2 - If the Idea's Thought is a map[string]* (such as a Nexus), this will Stringify the first path component and recursively navigate relatively inwards.
//
// 3 - Otherwise, this will recursively reflect into the Idea's Revelation() method using the below pathing rules:
//
//	# 0 Fields are simply accessed
//	# 1 Methods will be called, but can only take a std.Impulse
//	# 2 The following method signatures can be implicitly called:
//		func()   | func(*)   | func(...*)   | func(*, ...*)   -> yields a std.Void structure
//		func() * | func(*) * | func(...*) * | func(*, ...*) * -> yields the returned wildcard
//	# 3 The current std.Impulse will only be injected automatically into the first parameter
func (id *Idea[t]) Relative(path any, code ...any) (any, error) {
	// id.Relative("") -> Revelation()
	// id.Relative("Revelation") -> Revelation().Revelation___
	// id.Relative("Field") -> Revelation().Field___
	// id.Relative("Method1", "Field", "Method2") -> Revelation().Method1().Field.Method2()
}
