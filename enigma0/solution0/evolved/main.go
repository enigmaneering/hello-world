package main

import (
	"fmt"

	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/enum/lifecycle"
	"git.ignitelabs.net/janos/core/std"
	"git.ignitelabs.net/janos/core/sys/atlas"
	"git.ignitelabs.net/janos/core/sys/when"
)

func init() {
	atlas.Verbose(true)
}

/*
E0S0

This prints several statistics through phasing activation across time using two different techniques:

	- Cortical Phasing, where the impulse.Beat value creates a manual round-robin effect
	- Clustered Phasing, where a neural cluster implicitly creates a stable and safe round-robin effect

The reason cortical phasing is unstable is that neurons can decay, causing delayed gaps in the impulse cycle
which don't get "refilled."  When using a cluster, this issue is resolved and decayed neurons naturally fall
out of the cycle.  You can explicitly see this in the outputs of both:

In 'cortical phasing' the cycle period ALWAYS remains 1.5s, even after one neuron decays out.  While watching
the output, you'll see a 'stepped gait' about the printed output

	1.004188958s [refraction] 1.004188958s
	1.504204708s [response] 2.166µs
	1.499292125s [cycle] 1.499292125s
	[Alpha Fabler ⇝ Response Time Printer] decayed
	1.499211459s [response] 2.167µs
	1.500632417s [cycle] 1.500632417s
	1.500516666s [response] 750ns
	1.500013083s [cycle] 1.500013083s
	1.500342667s [response] 875ns

In 'clustered phasing' the cycle period drops naturally after the neuron decays out, beginning on the cluster's
next cycle.  While watching the output, the only 'hitch' is on the step where the neuron decays

	1.500055334s [response] 11.667µs
	1.50024425s [cycle] 1.50024425s
	1.499871792s [refraction] 1.499831625s
	[Ethelene Bottlestone ⇝ cluster ⇝ Response Time Printer] decayed
	1.50001375s [cycle] 1.50001375s
	1.499999s [refraction] 1.499959084s
	999.997083ms [cycle] 999.997083ms
	1.000005583s [refraction] 999.96475ms
	1.000013667s [cycle] 1.000013667s
*/

var cortex = std.NewCortex(std.RandomName())

func main() {
	cortex.Frequency = 2 //hz

	//UseCorticalPhasing()
	UseClusteredPhasing()

	cortex.Spark()
	core.KeepAlive()
}

func UseCorticalPhasing() {
	cortex.Phase = 2

	cortex.Synapses() <- std.NewSynapse(lifecycle.Looping, "Cycle Printer", PrintCycle, func(imp *std.Impulse) bool {
		return imp.Beat == 0
	})

	cortex.Synapses() <- std.NewSynapse(lifecycle.Looping, "Refractory Printer", PrintRefractory, func(imp *std.Impulse) bool {
		return imp.Beat == 1
	})

	decay := 0
	cortex.Synapses() <- std.NewSynapse(lifecycle.Looping, "Response Time Printer", PrintResponse, func(imp *std.Impulse) bool {
		if imp.Beat == 1 {
			decay++
			if decay == 2 {
				imp.Decay = true
			}
			return !imp.Decay
		}
		return false
	})
}

func UseClusteredPhasing() {
	// 0 - Create the synaptic cluster
	syn, cluster := std.NewSynapticCluster("cluster", cortex, when.Frequency(8))

	// 1 - Add it to the cortex
	cortex.Synapses() <- syn

	// 2 - Add neurons to the cluster
	cluster <- std.NewNeuron("Cycle Printer", PrintCycle, nil)
	cluster <- std.NewNeuron("Refractory Printer", PrintRefractory, nil)

	// 3 - Add a faux delayed decay to the final neuron for demonstration purposes
	decay := 0
	cluster <- std.NewNeuron("Response Time Printer", PrintResponse, func(imp *std.Impulse) bool {
		if decay == 2 {
			imp.Decay = true
		}
		decay++
		return !imp.Decay
	})
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
