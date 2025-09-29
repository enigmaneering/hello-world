package main

import (
	"fmt"

	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/enum/lifecycle"
	"git.ignitelabs.net/janos/core/std"
)

/*
E0S0

This prints several statistics through phasing activation across time
*/

func main() {
	c := std.NewCortex(std.RandomName())
	c.Frequency = 2 //hz
	c.Phase = 3

	c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Cycle Printer", PrintCycle, func(imp *std.Impulse) bool {
		return imp.Beat == 0
	})
	c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Refractory Printer", PrintRefractory, func(imp *std.Impulse) bool {
		return imp.Beat == 1
	})

	i := 0
	c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Response Time Printer", PrintResponse, func(imp *std.Impulse) bool {
		if imp.Beat == 2 {
			i++
			if i >= 3 {
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
}

func PrintRefractory(imp *std.Impulse) {
	fmt.Printf("%v [refraction] %v\n", imp.Timeline.CyclePeriod().String(), imp.Timeline.RefractoryPeriod().String())
}

func PrintResponse(imp *std.Impulse) {
	fmt.Printf("%v [response] %v\n", imp.Timeline.CyclePeriod().String(), imp.Timeline.ResponseTime().String())
}
