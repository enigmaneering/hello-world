package std

import (
	"fmt"
	"time"

	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/std"
)

type Synapse chan<- any

func NewSynapse(named string, action func(*Impulse), potential func(*Impulse) bool, cleanup ...func(*Impulse)) Synapse {
	signal := make(chan any, 1<<16) // NOTE: A 2¹⁶ is only a few kb of memory, but gives plenty of overhead for future access

	go func() {
		defer core.HandlePanic(named, "synaptic loop")

		act := Activation{
			Epoch: time.Now(),
		}
		imp := &Impulse{
			Entity:     std.NewEntityNamed(named),
			Activation: act,
			Bridge:     make(Bridge, 0),
			Timeline:   std.NewTemporalBuffer[Activation](),
		}

		panicSafeCleanup := func(i *Impulse) {
			defer core.HandlePanic(i.String(), "neural cleanup")

			if len(cleanup) > 0 && cleanup[0] != nil {
				cleanup[0](i)
			}
		}
		panicSafePotential := func(i *Impulse) bool {
			defer core.HandlePanic(i.String(), "neural potential")

			if potential == nil || potential(i) {
				return true
			}
			return false
		}
		panicSafeAction := func(i *Impulse) {
			defer core.HandlePanic(i.String(), "neural action")

			action(i)
		}

		for core.Alive() {
			msg := <-signal
			imp.Activation.Inception = time.Now()

			switch msg.(type) {
			case spark:
				if panicSafePotential(imp) {
					imp.Activation.Fired = core.Ref(time.Now())
					panicSafeAction(imp)
					imp.Activation.Completion = core.Ref(time.Now())
					imp.Timeline.Record(imp.Inception, imp.Activation)
				}
			default:
				panic(fmt.Errorf("unknown signal %T", msg))
			}
		}

		panicSafeCleanup(imp)
	}()

	return Synapse(signal)
}
