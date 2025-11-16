package main

import (
	"fmt"

	"git.ignitelabs.net/janos/core/std"
)

func main() {
	for i := 0; i < 10; i++ {
		fmt.Println(std.RandomName())
	}
}
