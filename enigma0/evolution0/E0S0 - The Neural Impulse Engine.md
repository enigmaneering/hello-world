# `E0S0 - The Neural Impulse Engine`
## `A.K.A. The Self-Regulating Looping Clock`
### `Alex Petz, Ignite Laboratories, March 2025`

---

### Jumping Gee Willikers - Looping Clocks!?

Well, yes!  And there's a fundamentally important reason for this:

_**Calculation takes time**_

Even the act of _reading_ a value takes time!  Subsequent layers of calculation only exacerbate
the issue, but intelligent systems _demand_ rich calculations!  At any given instant, billions of neurons are
firing across your biological machinery - accumulating thresholds, locking semaphores, and calculating _thought._

_How!?_

By harnessing the _shared_ passing of time through the power of the _feedback loop!_

### The Colonel of Kernels

The entirety of this project is built around a distributed operating system I've named `JanOS` after the Roman God of 
time, duality, gates, and beginnings.  Every Go program that imports the library becomes a JanOS **instance** by extension,
regardless of its integration into a larger system.  Your instance is given a random name at startup, though you can
set the defaults through the `atlas` configuration system (which we'll discuss later).  The instance provides global 
synchronization over any neural activity you'd like to **spark** off.

What does it mean to _spark_ something off?  Well, a goroutine can either execute a single operation or act as a long-running 
process.  When I refer to 'sparking' something, I'm specifically referring to the act of creating an _abstractly_ long-running
process.  Neural activity, while typically implemented through a loop, uses externally impulsed channels to gate the flow of 
logic.  Each of these little micro-kernels sustains several useful features:

0. The cyclic invocation of a neural action (abstractly, _your_ function)
1. The ability to gate the action's activation through potential functions
2. Panic-safe activation and cleanup
3. Runtime I/O - including signaling, temporal buffering, and activation stats

That's a lot to take in at once - and I haven't even gotten to the concept of phase-based activation and periodic
cycles!  Luckily, I've done my absolute best to make this API as clear-cut and intuitive as possible.  My model of a 
semantic-based neural operating system is _as fundamentally close_ to biological neurology as computer science 
currently affords - minus the chemicals, of course.  I don't aim to ensure that the result perfectly matches our 
biological machinery.  Instead, I hope to inspire you to consider the presence of a divine architect in the creation 
of our "source code" by putting _you_ in the hot-seat of creation!  Truly - my goal is to demonstrate the replication 
of a virtual "psyche" in code, able to animate _**its own**_ dream of a corporeal form as we co-author its reality 
into existence =)

That couldn't ever, _and would never,_ happen without _you!_  

...or Alice, Bob, and Charlie - foo, bar, and baz - or the absolutely _brilliant_ geniuses who dreamt those parallels before us!

I've written about a dozen different variants of this architecture now, and I'll probably write hundreds more in my
lifetime.  Each time, I've come to the realization that the _only_ way one could feasibly document such a beast
is by writing the documentation _concurrently._  To that end, I've chosen to _enigmaneer_ this round as a kind of stream
of consciousness.  Lucky for you, this _isn't_ the first draft!  You'll have to take into consideration that every
decision from here on out has an intelligent reason behind it, even if it feels highly anti-pattern at first.

Let's begin =)

### Primitive Puzzle Pieces

To begin with, let's define our flagship component - the synaptic neural channel!  What even _is_ a synapse?

    synapse
      noun — /ˈsɪnæps/
        1. a junction between two nerve cells, consisting of a minute gap 
           across which impulses pass by diffusion of a neurotransmitter.

Cutting through the hard physical constraints, I see a parallel to something in code: _the anonymous function!_
While an anonymous function, in of itself, isn't a synapse - it provides the _framework_ for how to construct a
synaptic activation that bridges abstract neural endpoints.  I'm not sure how to describe the pattern I've built
to model this design, but I'm calling it the 'Impulse Pattern'

    package std

    type Synapse chan<- any

We begin with the synapse - an anonymously typed channel.  Because this is anonymously typed, anyone could build
their own synapses that act on any signals they desire.  For JanOS, we won't expose the actual signal types -
instead, we'll expose a factory for _creating_ signals called a _signal maker_

    package std

    type signalMaker byte
    var Signal signalMaker

    type spark *Impulse

    func (signalMaker) Spark(source ...*Impulse) spark {
        if len(source) > 0 {
            return source[0]
        }
        return nil
    }

The first signal the maker will expose is the _Spark_ signal - which takes in a referential impulse.  The spark
is _unexported_ as you should always use the factory method for _creating_ signals.  My goal is to produce an API 
where the following...

    std.NewSynapse("Printer", func(*std.Impulse) { fmt.Println("Hello, World!") }, nil) <- std.Signal.Spark()

...will, in a single line, _impulse_ "Hello, World!" to the console.  Most importantly, I want it to have _zero_ additional
structures and configuration necessary to do this - instead, evolving the _signals_ for more and more advanced functionality.
That will make a lot more sense as we progress forward, but for now let's make the single-line impulse happen!

    package std

    import (
        "git.ignitelabs.net/janos/core"
        "git.ignitelabs.net/janos/core/std"
    )

    type Impulse struct { }

    type Synapse chan<- any
    
    func NewSynapse(named string, action func(*Impulse), potential func(*Impulse) bool, cleanup ...func(*Impulse)) Synapse {
	    signal := make(chan any, 1<<16) // NOTE: A 2¹⁶ is only a few kb of memory, but gives plenty of overhead for future access
    
        go func() {
            imp := &Impulse{
                Entity: std.NewEntityNamed(named),
            }

            for core.Alive() {
                msg := <-signal
    
                switch msg.(type) {
                case spark:
                    // NOTE: An absent potential means 'always activate' =)
                    if potential == nil || potential(imp) {
                        action(imp)
                    }
                default:
                    panic(fmt.Errorf("unknown signal %T", msg))
                }
            }

            // For now, just cleanup after system shutdown
            if len(cleanup) > 0 && cleanup[0] != nil {
                cleanup[0](imp)
            }
        }()
    
        return Synapse(signal)
    }

Here, we leverage the `std.Entity` type that JanOS's `core` package already provides.  A std.Entity gives a _unique
identifier_ and _name_ to any composed types and will get heavy use throughout this project.  In addition, we have
a `core.Alive()` method which globally indicates if the instance is currently alive.  All JanOS instances are _alive
on creation_ until a call to `core.Shutdown()` - which then sets alive to false, allowing any long-running loops to
naturally decay off of a shared signal.  Since JanOS is built off of _many_ isolated loops, this is a crucial value
to observe!  At this point, we've met the goal of our synaptic activation printing "Hello, World!" 

    func main() {
        std.NewSynapse("Printer", func(*std.Impulse) { fmt.Println("Hello, World!") }, nil) <- std.Signal.Spark()
    }

...but we still have several of the features I mentioned at the start of this document to implement.  First, we'll 
start with _panic safe activation_ by wrapping the function calls

    package core

    // HandlePanic safely recovers from a panic, prints the panic event out, and optionally includes the debug.Stack().
    //
    // NOTE: If verbose is omitted, this will follow atlas.Verbose()
    func HandlePanic(named string, location string, verbose ...bool) {
        v := atlas.Verbose()
        if len(verbose) > 0 {
            v = verbose[0]
        }
    
        if r := recover(); r != nil {
            if v {
                fmt.Printf("[%s] %s panic: %v\n%s", named, location, r, debug.Stack())
            } else {
                fmt.Printf("[%s] %s panic: %v\n", named, location, r)
            }
        }
    }

    package std
    
    func NewSynapse(named string, action func(*Impulse), potential func(*Impulse) bool, cleanup ...func(*Impulse)) Synapse {
	    signal := make(chan any, 1<<16) // NOTE: A 2¹⁶ is only a few kb of memory, but gives plenty of overhead for future access
    
        go func() {
            defer core.HandlePanic(named, "synaptic loop")
    
            imp := &Impulse{
                Entity: std.NewEntityNamed(named),
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
    
                switch msg.(type) {
                case spark:
                    if panicSafePotential(imp) {
                        panicSafeAction(imp)
                    }
                default:
                    panic(fmt.Errorf("unknown signal %T", msg))
                }
            }
    
            panicSafeCleanup(imp)
        }()
    
        return Synapse(signal)
    }

Here, we've offloaded panic recovery to `core.HandlePanic`, which has some unique quirks.  First, let's talk about
`atlas` and verbosity.  Every JanOS instance can be _configured_ using the atlas configuration system.  This consists of
including an `atlas` file with your `go.mod`, or by setting the values at runtime (often through an `init()` function).
Verbosity, as you would imagine, indicates if the instance should emit more detailed logging values.  If you'd
like to debug a neural panic, pass _true_ to the optional final parameter of `HandlePanic` to override the default verbosity
settings.  Otherwise, neural panics leave no stack information in the output.  The `HandlePanic` function lives in the `core`
module, which is also where the `std` package we're currently adding to _also_ lives.  If you'd like to set the verbosity globally, 
the easiest way is through an init function (though you can change it at any time)

    func init() {  
        // Gets and/or sets the global verbosity value

        // Setting
        atlas.Verbose(true)

        // Getting
        if atlas.Verbose() {
            ...
        }
    }

Let's review our desired feature set -

0. The cyclic invocation of a neural action (abstractly, _your_ function)
1. The ability to gate the action's activation through potential functions
2. Panic-safe activation and cleanup
3. Runtime I/O - including signaling, temporal buffering, and activation stats

Out of these we are missing cyclic activation, temporal buffering, and activation stats (along with more robust
control signals).  Temporal buffering is accomplished through the `std.TemporalBuffer` type, which provides
a few ways to record and grab a rolling window of observed temporal data.  Temporal buffers are provided an _observation
window_ value indicating how long to hold data before decaying it out - if omitted, it defaults to `atlas.ObservanceWindow`

While the inner workings aren't important at the moment (it's built like any other buffer), it _does_ have a few
important quirks when working with it.  First, _all temporal data is zero indexed to epoch_ - which is just fancy speak for
_"index 0 comes temporally before index 𝑛."_  Second, these are the most critical methods it exposes

    var tbf TemporalBuffer[T any]

    tbf.Yield() // Retrieves all current elements
    tbf.Record(time.Time, T) // Records an element of type T with the provided moment of time

    tbf.Latest(...𝑛) // Retrieves up to the latest 𝑛-depth of elements, or a single entry if 𝑛 is omitted
    tbf.LatestSince(time.Time, ...bool) // Retrieves the latest elements since the provided moment, optionally including it

`LatestSince` will get used _heavily_ as we progress forward, but I'd like to briefly note _why_ you'd optionally
include the requested moment.  One of the features of JanOS is _temporal analysis_ using calculus.  This happens by
cyclically differentiating or integrating on a rolling range of data moments.  Sometimes you'll want to calculate only on
the distinctly new moments (exclusive), while other times you'll need to calculate the differential _between_ momentary data
(inclusive).  One of the tenants of JanOS is intuitively providing _options,_ thus "exclusive" is the default mode of operation 
before more advanced use cases like differential analysis of data.

Philosophically, that resonates as well!  Our default mode of operation is _exclusive_ of others, where inclusivity is
what we _"opt-in"_ to.  If I chose the opposite default condition, I'd imply that we selectively opt into exclusivity - but 
Life is a _feedback loop!_  If you don't _**dance**_ with it, you'll find yourself in a stagnant state of operation.  Holding 
relationships with others is a _choice,_ not a privilege - and accepting that perspective is quite empowering!  It injects 
_intention_ into your decisions, rather than apathy =)

Anywho!  Let's build some runtime stats into our synapse

    type Activation struct {
        id uint64
    
        // Epoch represents the moment the synaptic connection was created.
        Epoch time.Time
    
        // Inception represents the moment a neuron received a synaptic event.
        Inception time.Time
    
        // Fired represents the moment a neuron passed it's potential.
        Fired *time.Time
    
        // Completion represents the moment a neuron finished execution.
        Completion *time.Time
    }
    
    func (a *Activation) ID() uint64 {
        return a.id
    }

The ID is protected for our benefit, as it's used to find the activation in the temporal buffer for updates and
should never get set after creation.  An ounce of prevention is worth a pound of cure, even from ourselves!  The
`Fired` and `Completion` times are pointer types as they are _nil_ until they have a settable value.  Next,
we stubbed out the `Impulse` type at the start - let's flesh it out

    type Impulse struct {
        std.Entity // Identification
	    Bridge // The chain that created this impulse
        Activation // Runtime statistics
	    Synapse // The signalling channel
    
        Timeline *std.TemporalBuffer[Activation]
    }

    func (imp Impulse) String() string {
        if len(imp.Bridge) == 0 {
            return imp.Name.Name
        }
        return imp.Bridge.String() + " ⇝ " + imp.Name.Name
    }

    type Bridge []string
    
    func (b Bridge) String() string {
        return strings.Join(b, " ⇝ ")
    }

Right now, this type is pretty straightforward - it just composes some crucial context for the neural endpoint to leverage.
The bridge is used to identify a chain of activations.  For now, the chain of activation is just the synapse but,
as we build deeper layers of neural activity, this will obviously grow.  The bridge finds its use through JanOS's
built-in logging package, `rec` - which prints the caller's name in square brackets before otherwise acting as `fmt.Printf`

    rec.Printf(imp.String(), "%v!", "Hello, World") // [impulse name] Hello, World!

Now, let's add some methods to our impulse structure for calculating some important activation stats

    // RefractoryPeriod represents the duration between the last activation's completion and this impulse's firing moment.
    func (imp *Impulse) RefractoryPeriod() time.Duration {
        last := imp.Timeline.Latest()
        if len(last) <= 0 {
            return 0
        }
        return imp.Activation.Fired.Sub(*last[0].Element.Completion)
    }
    
    // CyclePeriod represents the duration between the last and current activation firing moments.
    func (imp *Impulse) CyclePeriod() time.Duration {
        last := imp.Timeline.Latest()
        if len(last) <= 0 {
            return 0
        }
        return imp.Activation.Fired.Sub(*last[0].Element.Fired)
    }
    
    // ResponseTime represents the duration between inception and firing of the current impulse.
    func (imp *Impulse) ResponseTime() time.Duration {
        return imp.Activation.Fired.Sub(imp.Activation.Inception)
    }
    
    // RunTime represents the duration between firing and completion of the last impulse's activation.
    func (imp *Impulse) RunTime() time.Duration {
        last := imp.Timeline.Latest()
        if len(last) <= 0 {
            return 0
        }
        return last[0].Element.Completion.Sub(*last[0].Element.Fired)
    }
    
    // TotalTime represents the duration between inception and completion of the last impulse's activation.
    func (imp *Impulse) TotalTime() time.Duration {
        last := imp.Timeline.Latest()
        if len(last) <= 0 {
            return 0
        }
        return last[0].Element.Completion.Sub(last[0].Element.Inception)
    }

Next, we get to populate the activation stats in our synaptic loop

    act := Activation{
        Epoch: time.Now(),
    }
    imp := &Impulse{
        Entity:     std.NewEntityNamed(named),
        Activation: act,
        Bridge: make(Bridge, 0),
        Timeline:   std.NewTemporalBuffer[Activation](),
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

To put it all together, let's output the response time of our primitive synapse!

    func main() {
        std.NewSynapse("Printer", func(imp *std.Impulse) {
            rec.Printf(imp.String(), "Hello, World! (%v)\n", imp.ResponseTime())
        }, nil) <- std.Signal.Spark()
    }
---
       ╭─⇥ JanOS v0.0.42
    ╭──┼───────⇥ © 2025 - Humanity
    ⎨  ⎬⇥ Maintained by The Enigmaneering Guild
    ╰─┬┴───⇥ ↯ [core] "Leonie Luckwell" - Lion
      ╰──────⬎
    [Printer] Hello, World! (333ns)

Wonderful!

Congratulations - you've just built the most basic scaffolding from which to build neural architectures!
I'd be remiss if I forced my system's preamble upon you, though.  The preamble can _only_ be disabled through an `atlas` file
placed alongside your `go.mod` - currently, a simple JSON format
    
    {"printPreamble": false}

### Refractory Periods

One of the most interesting statistics to consider is the synaptic _refractory period!_  This represents the span of time
between the _end_ of an activation and the _firing moment_ of the next - which provides a diminishing value as your action grows
in runtime.  This will come in quite handy later on as we explore distributing calculation against a window of time!  Let's
polish up our first solution by printing out the refractory period in a loop

    func main() {
        syn := std.NewSynapse("Printer", func(imp *std.Impulse) {
            rec.Printf(imp.String(), "Hello, World! (%v)\n", imp.RefractoryPeriod())
        }, nil)
    
        for core.Alive() {
            syn <- std.Signal.Spark()
    
            time.Sleep(time.Second)
        }
    }
---
    [Printer] Hello, World! (0s)
    [Printer] Hello, World! (1.001123375s)
    [Printer] Hello, World! (1.000595125s)
    [Printer] Hello, World! (1.001005625s)

In the next solution, we'll start evolving our architecture to sustain thoughts =)