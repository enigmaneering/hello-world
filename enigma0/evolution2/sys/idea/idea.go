package idea

import (
	"fmt"

	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution2/std"
)

type path[T any] std.Path

func Path[T any](path ...fmt.Stringer) path[T] {
	return path
}

func (p path[T]) Get(password ...string) T {

}

func (p path[T]) Set(value T, password ...string) {

}
