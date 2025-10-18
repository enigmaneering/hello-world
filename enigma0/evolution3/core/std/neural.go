package std

import "sync"

type Neural[T any] interface {
	Named() string

	Action(Impulse)

	Potential(Impulse) bool

	Cleanup(Impulse, *sync.WaitGroup)

	Heart() *Heart
}
