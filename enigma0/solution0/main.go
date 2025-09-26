package main

import (
	"fmt"

	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/enum/lifecycle"
	"git.ignitelabs.net/janos/core/std"
	"git.ignitelabs.net/janos/core/sys/when"
)

/*
E0S0

This prints several statistics through a stepped potential function.
*/

var Step func()

func main() {
	c := std.NewCortex(std.RandomName())
	c.Frequency = 60 //hz

	makePotential, step := when.StepMaker(3) // There are 3 neural endpoints to step between
	Step = step

	c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Cycle Printer", PrintCycle, makePotential(0))
	c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Refractory Printer", PrintRefractory, makePotential(1))

	i := 0
	final := makePotential(2)
	c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Response Time Printer", PrintResponse, func(imp *std.Impulse) bool {
		if final(imp) {
			i++
			if i >= 128 {
				i = 0
				return true
			}
		}
		return false
	})

	c.Spark()
	core.KeepAlive()
}

func PrintCycle(imp *std.Impulse) {
	fmt.Printf("%v [cycle] %v\n", imp.Timeline.CyclePeriod().String(), imp.Timeline.CyclePeriod().String())
	Step()
}

func PrintRefractory(imp *std.Impulse) {
	fmt.Printf("%v [refraction] %v\n", imp.Timeline.CyclePeriod().String(), imp.Timeline.RefractoryPeriod().String())
	Step()
}

func PrintResponse(imp *std.Impulse) {
	fmt.Printf("%v [response] %v\n", imp.Timeline.CyclePeriod().String(), imp.Timeline.ResponseTime().String())
	Step()
}
