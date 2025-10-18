package main

import (
	"fmt"
	"time"

	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/sys/id"
)

func main() {
	for core.Alive() {
		fmt.Println(id.Next())
		time.Sleep(time.Second)
	}
}
