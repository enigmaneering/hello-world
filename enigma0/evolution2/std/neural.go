package std

import "sync"

type Neural[T any] interface {
	Named() string

	Action(Impulse[T])

	Potential(Impulse[T]) bool

	Cleanup(Impulse[T], *sync.WaitGroup)
}
