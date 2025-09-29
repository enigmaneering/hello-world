package main

import (
	"time"

	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/enum/lifecycle"
	"git.ignitelabs.net/janos/core/std"
	"git.ignitelabs.net/janos/core/sys/rec"
)

/*
E0S2

This demonstrates how to clean up a neural endpoint
*/

func main() {
	c := std.NewCortex(std.RandomName())
	c.Frequency = 1 //hz

	c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Print", Printer, nil, Cleanup)

	go core.Shutdown(time.Second * 3)

	c.Spark()
	core.KeepAlive(time.Second * 5)
}

func Cleanup(imp *std.Impulse) {
	rec.Printf(imp.Bridge.String(), "cleaning up\n")
}

func Printer(imp *std.Impulse) {
	rec.Printf(imp.Bridge.String(), "%v\n", imp.Timeline.CyclePeriod())
}
