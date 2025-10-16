package std

import (
	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution2/core/enum/errs"
	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution2/core/enum/relationship"
)

// An Idea represents a locatable and psychologically constrainable Thought.
//
// NOTE: Please use the idea.New, idea.Reveal, and idea.Describe functions to abstractly access the idea std.Nexus.
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

func (id *Idea[T]) getThought(code ...any) (*Thought[T], error) {
	if id.disclosure.Constraint == relationship.Open {
		return id.thought, nil
	}

	if len(code) == 0 {
		return nil, errs.IdeaCodeRequired
	}
	if Stringify(code[0]) != Stringify(id.disclosure.Code) {
		return nil, errs.IdeaCodeInvalid
	}
	return id.thought, nil
}

func (id *Idea[T]) Reveal(code ...any) (T, error) {
	t, err := id.getThought(code...)
	if err == nil {
		return t.revelation, nil
	}
	var zero T
	return zero, err
}

func (id *Idea[T]) Describe(revelation T, code ...any) error {
	t, err := id.getThought(code...)
	if err == nil {
		t.Revelation(revelation)
		return nil
	}
	return err
}
