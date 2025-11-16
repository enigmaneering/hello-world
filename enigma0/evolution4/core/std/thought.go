package std

import (
	"sync"

	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution4/core/enum/relationally"
)

// A Thought is a thread-safe and relationally.Constrained revelation.  When used as a LIQ, the inner
// revelation is the Stringable component unless set by Thought.Stringable.
type Thought[T any] struct {
	revelation T
	gate       *sync.Mutex
	disclosure *Disclosure
	stringable func() string
	created    bool
}

func NewThought[T any](revelation T, disclosure ...*Disclosure) (*Thought[T], *Disclosure) {
	d := &Disclosure{
		Constraint: relationally.Open,
		code:       nil,
	}
	if len(disclosure) > 0 && disclosure[0] != nil {
		d = disclosure[0]
	}

	return &Thought[T]{
		revelation: revelation,
		gate:       new(sync.Mutex),
		disclosure: d,
		created:    true,
	}, d
}

func (t *Thought[T]) sanityCheck() {
	if !t.created {
		panic("std.Thought must be created through std.NewThought[T]")
	}
}

func (t Thought[T]) String() string {
	if t.stringable != nil {
		return t.stringable()
	}
	return Stringify(t.revelation)
}

// Stringify sets the method for stringing this Thought out.  If nil (default), Stringable operates directly on the
// revelation - otherwise, the provided string function is called.
func (t *Thought[T]) Stringify(fn func() string) {
	t.stringable = fn
}

// Reveal returns the underlying revelation of this Thought.
//
// NOTE: To reveal a relative path, please use Relative.
func (t *Thought[T]) Reveal(code ...any) (T, error) {

}

// Describe sets the underlying revelation of this Thought.
func (t *Thought[T]) Describe(revelation T, code ...any) error {

}

// Relative walks the provided Path and yields its target.
//
// NOTE: The code is used at any constrained points in the path, otherwise it's ignored. If you
// need to use multiple codes, you must sequentially reveal each codified part of the path.
func (t *Thought[T]) Relative(relative Path, code ...any) (any, error) {
	current, err := t.Reveal(code...)
	if err != nil {
		return nil, err
	}
	if len(relative) > 1 {
		relative = relative[1:]
		current, err =
	}
	return current, nil
}
