package std

import (
	"fmt"
	"time"

	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/std"
)

type Synapse chan<- any

func NewSynapse(named string, action func(*Impulse), potential func(*Impulse) bool, cleanup ...func(*Impulse)) Synapse {
	signal := make(chan any)

	go func() {
		defer core.HandlePanic(named, "synaptic loop")

		panicSafeCleanup := func(imp *Impulse) {
			defer core.HandlePanic(named, "neural cleanup")

			if len(cleanup) > 0 && cleanup[0] != nil {
				cleanup[0](imp)
			}
		}
		panicSafePotential := func(imp *Impulse) bool {
			defer core.HandlePanic(named, "neural potential")

			if potential == nil || potential(imp) {
				return true
			}
			return false
		}
		panicSafeAction := func(imp *Impulse) {
			defer core.HandlePanic(named, "neural action")

			action(imp)
		}

		act := Activation{
			Creation: time.Now(),
		}
		imp := &Impulse{
			Entity:     std.NewEntityNamed(named),
			Synapse:    signal,
			Activation: act,
			Timeline:   std.NewTemporalBuffer[Activation](),
		}
		decayCount := decay(0)
		decayClose := false
		decaying := false

		for core.Alive() {
			raw := <-signal
			if decaying && decayCount == 0 {
				break
			}

			imp.Activation.Inception = time.Now()

			switch msg := raw.(type) {
			case spark:
				if panicSafePotential(imp) {
					imp.Activation.Activation = core.Ref(time.Now())
					panicSafeAction(imp)
					imp.Activation.Completion = core.Ref(time.Now())
					imp.Timeline.Record(imp.Inception, imp.Activation)

					if decayCount > 0 {
						decayCount--
					}
					if decayCount < 0 {
						decayCount++
					}
				}
			case decay:
				decayCount = msg
				if decayCount < 0 {
					decayClose = true
				}
				decaying = true
			case nil:
			default:
				panic(fmt.Errorf("unknown signal %T", msg))
			}
		}

		go panicSafeCleanup(imp)

		// Enter silent mode
		for core.Alive() {
			if decayClose {
				close(signal)
				return
			}

			raw := <-signal

			if _, ok := raw.(closer); ok {
				close(signal)
				return
			}
		}
	}()

	return Synapse(signal)
}
