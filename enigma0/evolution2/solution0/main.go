package main

import (
	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/sys/atlas"
	"git.ignitelabs.net/janos/core/sys/rec"
	"time"
)

func init() {
	atlas.Verbose(true)
	atlas.Silent(true)
	atlas.ObservanceWindow = time.Second * 5
}

func main() {
	syn := std.NewSynapse("Printer", func(imp *std.Impulse[bool]) {
		rec.Printf(imp.String(), "Hello, World! (%v)\n", imp.RefractoryPeriod())
	}, nil)

	for core.Alive() {
		syn <- std.Signal[bool]().Spark(std.NewThought(true))

		time.Sleep(time.Second)
	}
}
