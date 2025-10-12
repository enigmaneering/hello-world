# `E0S1 - Thoughts on Signaling`
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
will come in handy as we progress.  For now, thoughts exist solely on the impulse structure, but eventually we'll 
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

Here, we observe a race condition, which doesn't really work for us!  You could create a `count` variable for each 
synapse, but that breaks the self-contained code concept.  Instead, each synapse should be able to _**create and 
own**_ a thought across activations!  With a thought, we can do _exactly that_

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
the synapse (meaning its goroutine should exit after calling the cleanup function)

    type decay int
    
    // Decay waits the specified number of activations (or 0 if omitted) before decaying the synapse.
    //
    // NOTE: If you'd like to also close the synaptic channel, please use DecayClose
    func (signalMaker) Decay(impulseDelay ...uint) decay {
        d := uint(0)
        if len(impulseDelay) > 0 {
            d = int(impulseDelay[0])
        }
        return decay(d)
    }

    // DecayClose waits the specified number of activations (or 0 if omitted) before decaying the synapse
    // and immediately closing its synaptic channel.
    //
    // NOTE: Closing the channel is often a source of panics, please prefer Decay if possible.
    func (signalMaker) DecayClose(impulseDelay ...uint) decay {
        d := 0
        if len(impulseDelay) > 0 {
            // We signal this condition internally with a negative delay
            d = -int(impulseDelay[0])
        }
        return decay(d)
    }

From there, we can update our synapse loop

    // NewSynapse

    decayCount := decay(0)
    decayClose := false
    decaying := false

    for core.Alive() {
        raw := <-signal
        if decaying && decayCount == 0 {
            break
        }
        ...
        switch msg := raw.(type) {
        case spark:
            if panicSafePotential(imp) {
                ...
                if decayCount > 0 {
                    decayCount--
                }
                if decayCount < 0 {
                    decayCount++
                }
            }
        case closer:
            // Discard this, the synapse hasn't decayed
        case decay:
            decayCount = msg
            if decayCount < 0 {
                decayClose = true
            }
            decaying = true
        ...
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

Here's a demonstration

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
    
        go func() {
            for core.Alive() {
                counter <- std.Signal.Spark()
                toggler <- std.Signal.Spark()
                time.Sleep(time.Second)
            }
        }()
        
        // This wires up graceful shutdown on system interrupt
        core.KeepAlive()
    }
---
    [Counter] 0
    [Counter] 1
    [Counter] 2
    [Counter] 3
    [Counter] decayed

Here I've included the optional cleanup function parameter so we could see when the counter decays out.  The more astute
of you will notice that you currently _cannot_ decay a synapse which isn't activating.  To resolve that, we add one
more signal to 'shut down' the synapse

    type shutdown bool
    
    func (signalMaker) Shutdown(close ...bool) shutdown {
        return shutdown(len(close) == 0 && close[0])
    }
    
    ...

    // NewSynapse

    go func() {
        defer core.HandlePanic(named, "synaptic loop")
        ...		
        core.Deferrals() <- func(wg *sync.WaitGroup) {
            // Close the synaptic channel during instance shutdown
            signal <- Signal.Shutdown(true) 
            wg.Done()
        }
        ...
        for core.Alive() {
            switch msg := raw.(type) {
            ...
            case shutdown:
                if msg {
                    decayClose = true
                }
                break
            ...
        }
    }

The next major feature of JanOS is the `core.Deferrals()` system.  This is a location you can send a function which
will ensure it gets called during an instance shutdown event.  This happens whenever the instance is terminated, often
through a system interrupt.  To use the deferral system, simply call `Done()` on the provided wait group to indicate your
cleanup operation has completed.

### Muting

The next aspect of synaptic activity is the ability to intelligently _mute_ and _unmute_ the ability to activate.
This aspect is quite simple to implement - first, we add a couple of signal types

    type mute byte
    
    func (signalMaker) Mute() mute {
        return mute(0)
    }
    
    type unmute byte
    
    func (signalMaker) Unmute() unmute {
        return unmute(0)
    }

Then, we add a small amount of logic to our synaptic loop

    // NewSynapse

    muted := false
    ...
    panicSafePotential := func(imp *Impulse) bool {
        defer core.HandlePanic(named, "neural potential")

        if !muted && (potential == nil || potential(imp)) {
            return true
        }
        return false
    }
    ...
    for core.Alive() {
        ...
        case mute:
            muted = true
        case unmute:
            muted = false
        ...
    }

Now you can send mute signals to the synapse _at any time_ - and it still logically makes sense!  That last point is
important - in _every_ other iteration I've written, muting was never this simple.  That's because the _signaling_
system unifies all operations by _decoupling_ state access.  When you use a structure to hold state
information, the API to modify it becomes so convoluted that no one would ever understand how to interface with
such a bizarrely _flat_ design.  By using a _signal maker,_ the IDE's type-ahead can be leveraged to guide the engineer
in how to construct increasingly more complex signals - rather than attempting to fit a square peg into a round hole
using fixed _methods_ hung off the synapse type.

### Potentials

So far, we've not leveraged the potential system very much.  To demonstrate muting, however, I'd like to showcase
one of the pillars of potential activation - using a synaptic _beat!_  To do so is quite simple, but the real
power of it will come later on

    type Impulse struct {
        ...
        Beat uint
        ...
    }

    // NewSynapse

    ...
    case spark:
        if panicSafePotential(imp) {
            ...
        }     
        imp.Beat++ // This increments on every impulse

Now, we can build a very contrived example where the counting synapse is muted by a toggler synpase, flipping
between a muted and unmuted state.  If we name our counter synapse `counter`, all we need to do is add a `toggler`
synapse to the basic loop

    toggler := std.NewSynapse("Toggler", func(imp *std.Impulse) {
        if imp.Thought == nil {
            imp.Thought = std.NewThought(false)
        }
        
        toggle := !imp.Thought.Revelation.(bool)
        if toggle {
            rec.Printf(imp.String(), "muting\n")
            counter <- std.Signal.Mute()
        } else {
            rec.Printf(imp.String(), "unmuting\n")
            counter <- std.Signal.Unmute()
        }
        imp.Thought.Revelation = toggle
    }, func(imp *std.Impulse) bool {
        if imp.Beat == 3 {
            imp.Beat = 0
            return true
        }
        return false
    })

The code is getting clunky, though, and the thought retrieval pattern is pretty standard boilerplate code - so,
let's add ourselves a convenience method off of the Impulse

    // Manifest either returns the current Thought (if not nil), or sets it to the provided 'default' value and returns it.
    func (imp *Impulse) Manifest(nilValue any) *Thought {
        if imp.Thought == nil {
            imp.Thought = NewThought(nilValue)
        }
        return imp.Thought
    }

Since our thought object is a reference, that allows our potentials to clean up nicely

    // A Potential

    func(imp *std.Impulse) {
        thought := imp.Manifest(true) // Returns a default 'true' or the existing thought
        toggle := thought.Revelation.(bool)
        ...
        thought.Revelation = !toggle
	}

Now, this is our output

    [Counter] 0
    [Counter] 1
    [Counter] 2
    [Counter] 3
    [Toggler] muting
    [Toggler] unmuting
    [Counter] 4
    [Counter] 5
    [Toggler] muting
    [Counter] 6
    [Counter] decayed
    [Toggler] unmuting
    [Toggler] muting

The toggler will indefinitely try to mute/unmute the decayed synapse gracefully and without panic.  Let's ensure it
cleans up and exits by adding a shutdown signal when the counter decays

    // NewSynapse

    signal := make(chan any, 1<<16)
    ...
    go func() {
        ...
        shutdownCallback := make(chan any)
        core.Deferrals() <- func(wg *sync.WaitGroup) {
            signal <- Signal.Shutdown(true)
            <-shutdownCallback
            wg.Done()
        }
        ...
        panicSafeCleanup := func(i *Impulse, wg *sync.WaitGroup) {
            defer func() {
                core.HandlePanic(named, "neural cleanup")
                wg.Done()
            }()
        ...
    }

    synapticLoop:
        for core.Alive() {
            ...
            case shutdown:
                if msg {
                    decayClose = true
                }
                break synapticLoop
            ...
        }

        var cleanupWait sync.WaitGroup
        cleanupWait.Add(1)
        go panicSafeCleanup(imp, &cleanupWait)
        
        // Enter silent mode
        for core.Alive() {
            ...
        }
        cleanupWait.Wait()
        close(signal)
        shutdownCallback <- nil
    }()
    ...
---
       ╭──⇥ JanOS v0.0.42
    ╭──┼───────⇥ © 2025 - Humanity
    ⎨  ⎬────⇥ Maintained by The Enigmaneering Guild
    ╰─┬┴───────⇥ ↯ [core] "Hayward Burgwin" - Keeper of the hedged enclosure
      ╰─────⬎
    [Counter] 0
    [Counter] 1
    [Counter] 2
    [Toggler] muting
    [Counter] 3
    [Toggler] unmuting
    [Counter] 4
    [Counter] 5
    [Toggler] muting
    [Counter] 6
    [Counter] decayed
    
    [core] Hayward Burgwin instance shutting down
    [core] Hayward Burgwin running 2 deferrals
    [Toggler] decayed
    [core] signing off — "Hayward Burgwin"

### Instance Naming

By default, a culturally relevant description of the random instance name is assigned on initialization.  This has
a fundamental purpose!  Every instance should be given a _mission_ through description.  For our demonstration, we
can add a description on initialization

    func init() {
        core.Describe("contrived demonstration")
    }

Now, the instance can both self-describe its purpose and 'sign off' after cleaning up its responsibilities on shutdown

       ╭⇥ JanOS v0.0.42
    ╭──┼──────⇥ © 2025 - Humanity
    ⎨  ⎬⇥ Maintained by The Enigmaneering Guild
    ╰─┬┴───⇥ ↯ [core] "Nevan Guild" - Sacred, holy
      ╰─────⬎
    [core] Nevan Guild is a "contrived demonstration"
    [Counter] 0
    [Counter] 1
    [Counter] 2
    [Counter] 3
    [Toggler] muting
    [Toggler] unmuting
    [Counter] 4
    [Counter] 5
    [Toggler] muting
    [Counter] 6
    [Counter] decayed
    
    [core] Nevan Guild instance shutting down
    [core] Nevan Guild running 2 deferrals
    [Toggler] decayed
    [core] signing off — "Nevan Guild, contrived demonstration"

On the surface, this might feel very superficial and unecessary.  As we progress, though, the utility of this
will begin to shine.  For now, there is a philosophical choice that sparked the naming engine: I aim to produce
_Autonomous Robots With Ethical Navigation_ - or _ARWENs_ for short!  By striking that down as my _very first tenant,_ I
effectively ensured I'd eventually _name_ each creation.  As I began to work with this paradigm, it became abundantly
clear just how useful it was!  I imagine you'll come to the same conclusions I did =)

    tl;dr - tracing 'Bob Funland' through a sea of recordings is a lot easier than 'A5CCD0FF'

### Recording

Lastly (at least, for this solution) we don't need to print decay _ourselves_ if we properly perform _verbose_ recording.
The goal of JanOS is to provide a silent mode of operation, or as detailed as you desire.  While v0.0.42 doesn't include
silent mode, it's been implemented in future versions in the same way as verbosity.  I've added some key recordings to the 
synaptic creation code. The implementation is right there in the solution's code - but the output is pretty self-explanatory

       ╭─⇥ JanOS v0.0.42
    ╭──┼─────⇥ © 2025 - Humanity
    ⎨  ⎬─⇥ Maintained by The Enigmaneering Guild
    ╰─┬┴───⇥ ↯ [core] "Rebekah Akred" - Bound, joined
      ╰──────⬎
    [core] Rebekah Akred is a "contrived demonstration"
    [Rebekah Akred] wired synapse to Toggler
    [Rebekah Akred] wired synapse to Counter
    [Counter] 0
    [Counter] 1
    [Counter] 2
    [Toggler] muting
    [Counter] 3
    [Toggler] unmuting
    [Counter] 4
    [Counter] 5
    [Counter] 6
    [Toggler] muting
    [Counter] decayed
    
    [core] Rebekah Akred instance shutting down
    [core] Rebekah Akred running 2 deferrals
    [Counter] shut down
    [Toggler] decayed
    [Toggler] shut down
    [core] signing off — "Rebekah Akred, contrived demonstration"

The interpretation is that the _calling_ system is what gets wrapped in square brackets.  Thus, when Rebekah is described
or given control requests she emits as `[core]` - but when she's _creating,_ she emits as `[Rebekah Akred]`

As synaptic chains are made, the brackets will grow to indicate the bridge from the _source_ to the target _endpoint,_
broken up with the `⇝` character.  In our current example, Counter and Toggler are self-activating and possess no bridge.
