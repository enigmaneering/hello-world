package std

import (
	"fmt"
	"sync"
	"time"

	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/std"
	"git.ignitelabs.net/janos/core/sys/rec"
)

type Synapse chan<- any

func NewSynapse(named string, action func(*Impulse), potential func(*Impulse) bool, cleanup ...func(*Impulse)) Synapse {
	signal := make(chan any, 1<<16)

	go func() {
		defer core.HandlePanic(named, "synaptic loop")

		shutdownCallback := make(chan any)
		core.Deferrals() <- func(wg *sync.WaitGroup) {
			signal <- Signal.Shutdown(true)
			<-shutdownCallback
			wg.Done()
		}

		act := Activation{
			Epoch: time.Now(),
		}
		imp := &Impulse{
			Entity:     std.NewEntityNamed(named),
			Synapse:    signal,
			Activation: act,
			Bridge:     make(Bridge, 0),
			Timeline:   std.NewTemporalBuffer[Activation](),
		}
		decayCount := decay(0)
		decayClose := false
		decaying := false
		muted := false

		panicSafeCleanup := func(i *Impulse, wg *sync.WaitGroup) {
			defer func() {
				core.HandlePanic(named, "neural cleanup")
				wg.Done()
			}()

			if len(cleanup) > 0 && cleanup[0] != nil {
				cleanup[0](i)
			}
		}
		panicSafePotential := func(i *Impulse) bool {
			defer core.HandlePanic(named, "neural potential")

			if !muted && (potential == nil || potential(i)) {
				return true
			}
			return false
		}
		panicSafeAction := func(i *Impulse) {
			defer core.HandlePanic(named, "neural action")

			action(i)
		}

		rec.Verbosef(core.Name.Name, "wired synapse to %v\n", imp.String())

	synapticLoop:
		for core.Alive() {
			raw := <-signal
			if decaying && decayCount == 0 {
				break
			}

			imp.Activation.Inception = time.Now()

			switch msg := raw.(type) {
			case spark:
				(*Impulse)(msg).BridgeTo(imp)
				if panicSafePotential(imp) {
					imp.Activation.Fired = core.Ref(time.Now())
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
				imp.Beat++
			case decay:
				decayCount = msg
				if decayCount < 0 {
					decayClose = true
				}
				decaying = true
			case mute:
				muted = true
			case unmute:
				muted = false
			case shutdown:
				if msg {
					decayClose = true
				}
				break synapticLoop
			case closer:
			case nil:
			default:
				panic(fmt.Errorf("unknown signal %T", msg))
			}
		}
		rec.Verbosef(imp.String(), "decayed\n")

		var cleanupWait sync.WaitGroup
		cleanupWait.Add(1)
		go panicSafeCleanup(imp, &cleanupWait)

		// Enter silent mode
		for core.Alive() {
			if decayClose {
				break
			}

			raw := <-signal

			if _, ok := raw.(closer); ok {
				break
			}
			if _, ok := raw.(shutdown); ok {
				break
			}
		}
		cleanupWait.Wait()
		close(signal)
		rec.Verbosef(imp.String(), "shut down\n")
		shutdownCallback <- nil
	}()

	return Synapse(signal)
}
