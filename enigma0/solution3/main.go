package main

import (
	"math/rand/v2"
	"time"

	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/enum/lifecycle"
	"git.ignitelabs.net/janos/core/std"
	"git.ignitelabs.net/janos/core/sys/rec"
)

/*
E0S3

This demonstrates the four different lifecycles of synaptic activity.
*/

var cortex = std.NewCortex(std.RandomName())

func main() {
	cortex.Frequency = 1 //hz

	// Choose your own adventure =)

	//Impulse()
	//Trigger()
	Stimulate()
	//Loop()
}

func Printer(imp *std.Impulse) {
	rec.Printf(imp.Bridge.String(), "START: %v\n", time.Now())
	time.Sleep(time.Second * 3)
	rec.Printf(imp.Bridge.String(), "STOP: %v\n", time.Now())
}

func Impulse() {
	/*
		An impulsive activation will decay after attempting to fire exactly once, regardless of success

		This example randomly generates a potential, then attempts to impulse the Printer function
	*/

	// Generate a random boolean value
	impulsePotential := rand.IntN(2) == 1
	cortex.Synapses() <- std.NewSynapse(lifecycle.Impulse, "Print", Printer, func(imp *std.Impulse) bool {
		if impulsePotential {
			rec.Printf(imp.Bridge.String(), "potential: high\n")
		} else {
			rec.Printf(imp.Bridge.String(), "potential: low\n")
		}
		return impulsePotential
	})

	cortex.Spark()
	core.KeepAlive()
}

func Trigger() {
	/*
		A triggered activation will fire exactly once while the potential is high, then decay

		This example sets up a trigger to activate the Printer function five seconds in the future
	*/

	now := time.Now().Add(time.Second * 5)
	cortex.Synapses() <- std.NewSynapse(lifecycle.Triggered, "Print", Printer, func(imp *std.Impulse) bool {
		if time.Now().After(now) {
			rec.Printf(imp.Bridge.String(), "potential: high\n")
			return true
		}
		rec.Printf(imp.Bridge.String(), "potential: low\n")
		return false
	})

	cortex.Spark()
	core.KeepAlive()
}

func Stimulate() {
	/*
		A stimulative lifecycle activates on every impulse the potential is high

		This example will activate the Printer function on every impulse indefinitely
	*/

	cortex.Synapses() <- std.NewSynapse(lifecycle.Stimulative, "Print", Printer, nil)

	cortex.Spark()
	core.KeepAlive(time.Second * 5)
}

func Loop() {
	/*
		A looping lifecycle recycles activation after each cycle completes execution

		This example will cyclically activate the Printer function indefinitely
	*/

	cortex.Synapses() <- std.NewSynapse(lifecycle.Looping, "Print", Printer, nil)

	cortex.Spark()
	core.KeepAlive(time.Second * 5)
}
