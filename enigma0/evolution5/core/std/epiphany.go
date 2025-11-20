package std

import (
	"sync"
	"time"

	"git.ignitelabs.net/janos/core"
)

// An Epiphany is a kind of Thought that exists in two states: in its idealized or materialized form.  For example, if you'd
// like to store data in JPEG format, its idealized form -is- encoded JPEG bytes - but its materialized form would decode
// those bytes into its raw color information.  An epiphany is a temporal structure, meaning it will materialize on request
// and remain as such before 'decaying' back into its idealized form after a period of inactivity by clearing the materialized
// result out of memory.
//
// Don't overthink it - that's really it =)
type Epiphany[TIdealized any, TMaterialized any] struct {
	thought    Thought[TIdealized]
	revelation *TMaterialized

	gate        *Gate
	motivate    chan any
	materialize func(TIdealized) (TMaterialized, error)
	decay       time.Duration
}

// NewEpiphany creates a new Epiphany which can 'materialize' into something more complex on demand, then 'decay' back to
// an idealized form after the provided amount of time with no activity.
func NewEpiphany[TIdealized, TMaterialized any](materialize func(TIdealized) (TMaterialized, error), decay time.Duration) *Epiphany[TIdealized, TMaterialized] {
	e := &Epiphany[TIdealized, TMaterialized]{
		motivate:    make(chan any),
		decay:       decay,
		materialize: materialize,
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

					// TODO: Cache the motivation as a "Decay" in std.Path{"Performance", "Epiphanies"}
				})

				// TODO: Cache the motivation as a "Reveal" in std.Path{"Performance", "Epiphanies"}
			}
		}
	}()

	return e
}

// StringifyFn sets the method for stringing this Epiphany out.  If nil (default), Stringable operates directly on the
// revelation - otherwise, the provided string function is called.
func (e *Epiphany[TIdealized, TMaterialized]) StringifyFn(fn func() string) {
	e.thought.StringifyFn(fn)
}

// Reveal returns the underlying revelation of this Epiphany.
//
// NOTE: To reveal a relative path, please use Recall.
func (e *Epiphany[TIdealized, TMaterialized]) Reveal(code ...any) (TMaterialized, error) {
	e.gate.Lock()
	defer e.gate.Unlock()

	select {
	case e.motivate <- nil:
	default:
		// No reason to block - contention is "good" for an epiphany
	}
	if e.revelation == nil {
		var ideal TIdealized
		var material TMaterialized
		var err error
		ideal, err = e.thought.Reveal(code)
		if err != nil {
			return material, err
		}

		material, err = e.materialize(ideal)
		if err != nil {
			return material, err
		}
		e.revelation = &material
	}
	return *e.revelation, nil
}

// Describe sets the underlying revelation of this Epiphany.
func (e *Epiphany[TIdealized, TMaterialized]) Describe(revelation TIdealized, code ...any) error {
	return e.thought.Describe(revelation, code...)
}

// Recall walks the provided Path relative to the current Epiphany and yields the result.
//
// NOTE: The code is used at any constrained points in the path, otherwise it's ignored. If you
// need to use multiple codes, you must sequentially reveal each codified part of the path.
func (e *Epiphany[TIdealized, TMaterialized]) Recall(relative Path) (any, error) {
	return e.thought.Recall(relative)
}
