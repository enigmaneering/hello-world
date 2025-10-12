package main

import (
	"time"

	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution1/std"
	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/sys/rec"
)

func main() {
	syn := std.NewSynapse("Printer", func(imp *std.Impulse[bool]) {
		rec.Printf(imp.String(), "Hello, World! (%v)\n", imp.RefractoryPeriod())
	}, nil)

	for core.Alive() {
		syn <- std.Signal[bool]().Spark(std.NewThought(true))

		time.Sleep(time.Second)
	}
}
