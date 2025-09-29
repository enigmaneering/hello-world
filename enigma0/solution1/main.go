package main

import (
	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/enum/lifecycle"
	"git.ignitelabs.net/janos/core/std"
	"git.ignitelabs.net/janos/core/sys/rec"
	"git.ignitelabs.net/janos/core/sys/when"
)

/*
E0S1

This toggles the activation frequency through the neural endpoint.
*/

func main() {
	c := std.NewCortex(std.RandomName())
	c.Frequency = 1 //hz

	rec.Printf(c.Named(), "Hello, World!\n")

	c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Print", Printer, when.FrequencyRef(&frequency))

	c.Spark()
	core.KeepAlive()
}

var frequency = 1.0
var toggle = true

func Printer(imp *std.Impulse) {
	// NOTE: This scheme doesn't guarantee an even timing interval =)

	if toggle {
		frequency = 0.5
	} else {
		frequency = 1.0
	}
	toggle = !toggle
	rec.Printf(imp.Bridge.String(), "%v\n", imp.Timeline.CyclePeriod())
}
