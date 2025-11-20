package std

import (
	"sync"
	"time"

	"git.ignitelabs.net/janos/core"
)

type Epiphany[TInner any, TOuter any] struct {
	thought    Thought[TInner]
	revelation *TOuter

	gate     *Gate
	motivate chan any
	evolveFn func(TInner) (TOuter, error)
	decay    time.Duration
}

func NewEpiphany[TInner, TOuter any](evolveFn func(TInner) (TOuter, error), decay time.Duration) *Epiphany[TInner, TOuter] {
	e := &Epiphany[TInner, TOuter]{
		motivate: make(chan any),
		decay:    decay,
		evolveFn: evolveFn,
	}

	go func() {
		var timer *time.Timer

		core.Deferrals() <- func(wg *sync.WaitGroup) {
			// On global shutdown, ensure the channel's loop is cleaned up
			e.motivate <- 1
			wg.Done()
		}
		for core.Alive() {
			select {
			case msg := <-e.motivate:
				if msg != nil {
					if timer != nil {
						timer.Stop()
					}

					return
				}
				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(e.decay, func() {
					e.gate.Lock()
					e.revelation = nil
					e.gate.Unlock()
				})
			}
		}
	}()

	return e
}

// StringifyFn sets the method for stringing this Epiphany out.  If nil (default), Stringable operates directly on the
// revelation - otherwise, the provided string function is called.
func (e *Epiphany[TInner, TOuter]) StringifyFn(fn func() string) {
	e.thought.StringifyFn(fn)
}

// Reveal returns the underlying revelation of this Epiphany.
//
// NOTE: To reveal a relative path, please use Recall.
func (e *Epiphany[TInner, TOuter]) Reveal(code ...any) (TOuter, error) {
	e.gate.Lock()
	defer e.gate.Unlock()

	select {
	case e.motivate <- nil:
	default:
		// No reason to block - contention is "good" for an epiphany
	}
	if e.revelation == nil {
		var inner TInner
		var outer TOuter
		var err error
		inner, err = e.thought.Reveal(code)
		if err != nil {
			return outer, err
		}

		outer, err = e.evolveFn(inner)
		if err != nil {
			return outer, err
		}
		e.revelation = &outer
	}
	return *e.revelation, nil
}

// Describe sets the underlying revelation of this Epiphany.
func (e *Epiphany[TInner, TOuter]) Describe(revelation TInner, code ...any) error {
	return e.thought.Describe(revelation, code...)
}

// Recall walks the provided Path relative to the current Epiphany and yields the result.
//
// NOTE: The code is used at any constrained points in the path, otherwise it's ignored. If you
// need to use multiple codes, you must sequentially reveal each codified part of the path.
func (e *Epiphany[TInner, TOuter]) Recall(relative Path) (any, error) {
	return e.thought.Recall(relative)
}
