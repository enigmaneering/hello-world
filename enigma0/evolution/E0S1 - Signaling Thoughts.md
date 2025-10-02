# `E0S0 - Signaling Thoughts`
### `Alex Petz, Ignite Laboratories, September 2025`

---

### Evolving Data

So far we've built a _very_ contrived way of printing the runtime of a single function call.  That's neat, but it'd
be a _lot_ more useful if the function could mature a thought over time!  So, what's a thought?

    type Thought struct {
        Revelation any
        Gate       *sync.Mutex
    }
    
    func NewThought(revelation any) *Thought {
        return &Thought{Revelation: revelation}
    }

That's it!  Your neuron can decide how complex the thought should be, but the ability to _lock_ access to the thought
will come in handy as we progress.  For now, thoughts exist solely on the impulse structure, but eventually you'll 
build a neural map from which thoughts can be shared between neurons

    type Impulse struct {
        std.Entity
        Bridge
        Activation
    
        Timeline *std.TemporalBuffer[Activation]
    
        Thought *Thought
    }

A thought is _only_ managed by the neural endpoint, making it a kind of inter-activation _cache._  Let's use it to
count!

    func main() {
        syn := std.NewSynapse("Counter", func(imp *std.Impulse) {
            if imp.Thought == nil {
                imp.Thought = std.NewThought(0)
            }
            
            count := imp.Thought.Revelation.(int)
            rec.Printf(imp.String(), "%v\n", count)
            imp.Thought.Revelation = count + 1
        }, nil)
    
        for core.Alive() {
            syn <- std.Signal.Spark()
    
            time.Sleep(time.Second)
        }
    }
---
    [Counter] 0
    [Counter] 1
    [Counter] 2
    [Counter] 3

The more astute of you would realize we could do this _a lot simpler_ with a closure!  However, a closure has a
_very important feature:_ it maintains its value across activations!  For example

    func main() {
        count := 0
        synA := std.NewSynapse("Counter A", func(imp *std.Impulse) {
            rec.Printf(imp.String(), "%v\n", count)
            count++
        }, nil)
        synB := std.NewSynapse("Counter B", func(imp *std.Impulse) {
            rec.Printf(imp.String(), "%v\n", count)
            count++
        }, nil)
    
        for core.Alive() {
            synA <- std.Signal.Spark()
            time.Sleep(time.Second)
            synB <- std.Signal.Spark()
        }
    }
---
    [Counter A] 0
    [Counter A] 1
    [Counter B] 2
    [Counter A] 3
    [Counter B] 3
    [Counter A] 5
    [Counter B] 5

Unfortunately, that doesn't really work for us!  Not only is the count split between activations, race conditions
are introduced from the closure.  Instead, each synapse should be able to _**create and own**_ a thought across 
activations!  With a thought, we can do _exactly that_

    func main() {
        createFn := func(named string) std.Synapse {
            return std.NewSynapse(named, func(imp *std.Impulse) {
                if imp.Thought == nil {
                    imp.Thought = std.NewThought(0)
                }
    
                count := imp.Thought.Revelation.(int)
                rec.Printf(imp.String(), "%v\n", count)
                imp.Thought.Revelation = count + 1
            }, nil)
        }
    
        synA := create("Counter A")
        synB := create("Counter B")
    
        for core.Alive() {
            synA <- std.Signal.Spark()
            synB <- std.Signal.Spark()
    
            time.Sleep(time.Second)
        }
    }
---
    [Counter A] 0
    [Counter B] 0
    [Counter B] 1
    [Counter A] 1
    [Counter B] 2
    [Counter A] 2
    [Counter B] 3
    [Counter A] 3

It's a subtle, but _powerful_ feature!  In addition, your synaptic revelation could be _anything_ - for instance, a 
`*TemporalBuffer[T]` would allow you to record temporal information on every impulse which decayed after a window of time.

### Signals

The next thing we get to work on is adding more signals to the synapse.  The counter above makes an excellent example
for the _decay_ signal!  A decay signal is a uint value which indicates the number of impulses to wait before decaying
the synapse (meaning it's goroutine should exit after calling the cleanup function).

    type decay int
    
    func (signalMaker) Decay(impulseDelay ...uint) decay {
        d := uint(0)
        if len(impulseDelay) > 0 {
            d = uint(impulseDelay[0])
        }
        return decay(d)
    }
    func (signalMaker) DecayClose(impulseDelay ...uint) decay {
        d := 0
        if len(impulseDelay) > 0 {
            d = -int(impulseDelay[0])
        }
        return decay(d)
    }

From there, we can update our synapse loop

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

At the point of decay, the synapse enters a 'silent mode' where all signals are discarded.  This is preferred over closing
the channel as Go doesn't provide any way to check if a channel is closed _before_ sending to it, which **_panics_**
if the channel is closed!  Instead, the graceful thing to do is to stop burning CPU cycles and let the goroutine
exit whenever the instance shuts down.  If you are certain it's safe to close the channel, you can send it a `Close`
signal - which will close the channel and exit the goroutine.

    type closer byte

    // Close will tell a decayed synapse it's safe to close its channel and fully exit.
    //
    // NOTE: If sent to a non-decayed channel, this will be discarded.
    func (signalMaker) Close() closer {
        return closer(0)
    }

And lastly, here's our demonstration

    func main() {
        syn := std.NewSynapse("Counter", func(imp *std.Impulse) {
            if imp.Thought == nil {
                imp.Thought = std.NewThought(0)
            }
    
            count := imp.Thought.Revelation.(int)
            rec.Printf(imp.String(), "%v\n", count)
            imp.Thought.Revelation = count + 1
        }, nil, func(imp *std.Impulse) {
            rec.Printf(imp.String(), "decayed\n")
        })
    
        syn <- std.Signal.Decay(4)
    
        for core.Alive() {
            syn <- std.Signal.Spark()
            time.Sleep(time.Second)
        }
    }
---
    [Counter] 0
    [Counter] 1
    [Counter] 2
    [Counter] 3
    [Counter] decayed

I've included the optional cleanup function parameter so we could see when the counter decays out.