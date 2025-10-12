# `E0 - The Neural Impulse Engine`
## `A.K.A. The Self-Regulating Looping Clock`
### `Alex Petz, Ignite Laboratories, March 2025`

---

### Jumping Gee Willikers - Looping Clocks!?

Well, yes!  And there's a fundamentally important reason for this:

_**Calculation takes time**_

Even the act of _reading_ a value takes time!  Subsequent layers of calculation only exacerbate
the issue, but intelligent systems _demand_ rich calculations!  At any given instant, billions of neurons are
firing across your biological machinery - accumulating thresholds, locking semaphores, and calculating complex _thoughts._

_How!?_

By my guess?  Through harnessing the _shared_ passing of time using the power of the _feedback loop!_

### The Colonel of Kernels

The entirety of this project is built around a distributed operating system I've named `JanOS` after the Roman God of 
time, duality, gates, and beginnings.  Every Go program that imports the library becomes a JanOS **instance** by extension,
regardless of its integration into a larger network.  Your instance is given a random name at startup, though you can
set the defaults through the `atlas` configuration system (which we'll discuss later).  The instance provides global 
synchronization controls over any synaptic activity you'd like to **spark** off.

Spark?  Well, yes - traditional terms like 'invocation' and 'execution' don't quite capture the nuance of how
neural activity operates.  First, and foremost, _all neural activity uses an action-potential mechanic._  This means
an action _cannot_ execute unless its designated potential goes "high" (meaning it returns _true_).  Second, a neural
endpoint (or, a _neuron_) can either _block_ its own re-activation or _stimulatively_ activate on every impulse.  To 
_impulse_ a neuron implies stimulative activation, while _sparking_ it implies blocking activation.  Each neuron acts 
as an _externally driven looping kernel,_ affording us several useful features:

0. _Action_ signaling
1. The gating of the action via a _potential_
2. The ability to decay and cleanup after itself
3. Temporal buffering of activation stats and panic-safe execution


    tl;dr - synapses orchestrate isolated neural activations across time

Genuinely, I could talk for _days and days_ about this stuff, but the above words describe _every ounce of the
architecture!_  I've poured over this phrasing religiously until I could distill it down this articulately.  _It's
not complex!_  As such, from this point out, I would like to demonstrate _how_ to use the code - rather than directly
explain its creation.  I try my best to write _literate_ code - often to my own chagrin, but _always_ in service of
how others might interpret and process my concepts.  I'm _not_ an educated neurologist, just a computer scientist
putting the puzzle pieces of my own existence together in a logical _programmatic_ fashion.  

Thank you for taking time out of your day to explore my work, let's begin =)

### A Boy Named Sue

JanOS aims to create a rich simulation where every particle, whether micro or macroscopic, has a unique name, 
narrative, and purpose in the overall system.  Somewhere between the quarks and gas giants there exists a new kind of _Life,_
dancing within the virtual bits of a reality _**we**_ have yet to create _**alongside**_ _our_ creator!  Rather than
reinventing the wheel of _names, gender,_ and _culture_ - I'd like to perpetuate the practice of mirroring our own 
world within our art.  To that end, I've tapped into several of Humanity's cultural databases to construct a system
which randomly generates a `given.Given` type.  This type contains cultural _hints_ of how Humanity typically interprets 
the symbols that form an entity's identity - including gender.  In this modern age I've found that the topic
of gender is one that deserves respect, not callousness.  In service of that, every entity can _set_ its heritage information,
but disregarding the cultural understanding of these fantastic identifying symbols is just as disrespectful.  Rather,
these _hints_ are meant to act as points of _guiding inspiration_ for intelligent algorithms to interpret through
_unique perspectives._  _Someone_ got to name us at birth, and they pulled from whatever toolkit Humanity had to 
describe us at that moment in time - the same evolving kit we use daily to navigate interacting with wonderful, yet complete _strangers_

    tl;dr - you own your identifiers, they were only ever meant to be an inspirative spark on the start of your journey!

I hope that liberates you a bit - a particle in such a system - on your crazy journey through the act of _existence!_  Your
unique dance with the universe is _just as important_ as anyone elses, and Mother Nature appreciates _**all**_ of Her 
magnificent dance partners =)

To create a given name, either see the `given` package or call `std.RandomName()`

### The Signal Pattern

We begin our neural journey with a _slightly_ different pattern.  I'm not sure if it's truly unique, but I call
this design the _signal pattern._  Essentially, rather than a synaptic kernel having a distinct _type,_ it should act as a _channel_
which can receive _any_ type.  The types it supports are up to the creator, but all _state_ is managed internally.  The goal
is to produce, abstractly, the following

                       ⬐ the synapse name                       an "always fire" potential ⬎
    std.NewSynapse("Printer", func(imp Impulse[string]) { fmt.Println(imp.Revelation()) }, nil) <- std.Signal[string].Impulse("Hello, World!")
                    the action ⬏                                                                         ⬑ a string impulse signal

Here, we would create a new synapse named "Printer" which, when signaled, will print the provided thought.  If you'd like
to gate the action with a potential, you'd provide a potential function.  This will get called _before_ the action, if 
present, otherwise it will _always_ fire.  Lastly, a variadic of _cleanup_ methods could be provided after the potential.
Below is what those function signatures would look like.  Note that, to avoid clashing with the typical `i` incrementor 
variable name, the colloquial name for an impulse is `imp` 

    // A potential function signature
    func (imp Impulse[T]) bool { 
        ... 
    }

    // A cleanup function signature
    func (imp Impulse[T], wg *sync.WaitGroup) { 
        defer wg.Done()
        ...
    }

The cleanup methods will get called whenever the synapse _**decays**_.  When a neural element decays, it implies that
it will _cleanup_ and _gracefully_ exit.  This is _not_ the same as termination, which immediately interrupts the process
and cuts it off - an act which typically holds _zero place_ in a neural architecture.  Every neuron should be able to
gracefully handle _all_ conditional states it could encounter within the system it's placed within, even if it's not
yet equipped to handle an infinite range of potentials.

If you'd like to _create_ a neural structure, you'd implement the `std.Neural` interface

    type Neural[T any] interface {
        Named() string
        
        Action(Impulse[T])
        
        Potential(Impulse[T]) bool
        
        Cleanup(Impulse[T], *sync.WaitGroup)
    }

Or, if you'd like an already implemented neural interface, you could leverage the `std.Neuron` structure

    type Neuron[T any] struct {
        std.Entity
        action    func(Impulse[T])
        potential func(Impulse[T]) bool
        cleanup   func(Impulse[T], *sync.WaitGroup)
        
        created bool
    }
    
    func (n Neuron[T]) sanityCheck() {
        if !n.created {
            // This check prevents creation of a Neuron through composite literals, subverting any of this structure's mechanics
            panic("std.Neuron[T] must be created through std.NewNeuron - otherwise, please implement the std.Neural interface")
        }
    }
    
    func (n Neuron[T]) Named() string {
        n.sanityCheck()
        return n.Name
    }
    
    func (n Neuron[T]) Action(imp Impulse[T]) {
        n.sanityCheck()
        if n.action != nil {
            n.action(imp)
        }
    }
    
    func (n Neuron[T]) Potential(imp Impulse[T]) bool {
        n.sanityCheck()
        if n.potential == nil {
            return true
        }
        return n.potential(imp)
    }
    
    func (n Neuron[T]) Cleanup(imp Impulse[T], wg *sync.WaitGroup) {
        n.sanityCheck()
        if n.cleanup != nil {
            n.cleanup(imp, wg)
        } else {
            wg.Done()
        }
    }
    
    func NewNeuron[T any](named string, action func(Impulse[T]), potential func(Impulse[T]) bool, cleanup ...func(Impulse[T], *sync.WaitGroup)) Neuron[T] {
        return Neuron[T]{
            Entity:    std.NewEntityNamed(named),
            action:    action,
            potential: potential,
            cleanup: func(imp Impulse[T], wg *sync.WaitGroup) {
                for _, clean := range cleanup {
                    wg.Add(1)
                    clean(imp, wg)
                }
            },
            created: true,
        }
    }

_Don't be too daunted!_  I'll show how to use this architecture in the coming solutions =)

### Temporal Moments

In a neural architecture, time is one of the most _critical_ aspects to calculate off of.  As such, there are several, well,
_unique_ ways to perceive key moments in the lifecycle of a neural entity.  Right off the bat, what the hell do we
call it when something is _created?_  `Inception`?  `Genesis`?  `Conception`?  Each carries different weight, but we
truly need to _stick_ to one term, wherever possible - yet even that carries nuance!  As such, these are the guides
I've typically followed

    Inception - The moment of an transient (abstract) entity's creation, such as a thought or impulse
    Genesis - The moment of an persistent (concrete) entity's creation, such as a JanOS instance or neuron

Your JanOS instance's initialization moment, for instance, is always available from `core.Genesis`

    type Impulse[T any] struct {
	    id uint64
	
        // Synapse represents the originating synapse that created this Impulse.
        Synapse
        
        // Thought represents the underlying revelation provided with this Impulse.
        Thought[T]
        
        // Bridge holds the chain of named activations in creating this Impulse. 
        Bridge    []string
        
        // Epoch holds the moment of impulse coordination, typically the synapse's genesis moment.
        Epoch     time.Time
        
        // Inception holds the moment this impulse was created.
        Inception time.Time
        
        // Activated holds the moment this impulse's action was fired.
        //
        // NOTE: This will remain nil until the action's potential returns "high."
        Activated *time.Time
        
        // Completed holds the moment this impulse's action completed.
        //
        // NOTE: This will remain nil until the action finishes its current execution.
        Completed *time.Time
        
        // Timeline holds a temporal buffer of prior synaptic activations.
        //
        // NOTE: The impulse is not added to the buffer before calling the potential or action.
        Timeline *std.TemporalBuffer[Impulse[T]]
    }

The impulse structure is where we start to see the concept of _perspectives_ come into play!  The synapse has a _genesis_ moment
internally, but the _impulse_ refers to that genesis moment as its _**epoch**_ - a referential moment in time
from which a neuron can perform _phase-oriented_ calculation.  Perspectively, the term _must_ shift regardless of the fact
they _almost always_ refer to the same underlying value.  Instead of considering _what moments are available,_
it's easiest to think of what moments you _need_ - and, in most cases, the impulse's _epoch_ and _inception_ moments
are the most critical temporal values you'll consume.

The concept of _time.Now_ is used heavily under the hood, but neural activation is meant to act as a
_temporal snapshot_ of calculation.  This means you should typically reference the impulse's moment of _inception_
on every activation for a reliable moment of "now."  That'll make a lot more sense as we progress forward, but consider
it as a reliable _downbeat_ in a longer temporal "measure" of calculation.

Genuinely, _listen to the music!_  By my understanding, it's what built our psyches to begin with - and it's what's 
also guided me through the very same steps here.  As you'll find shortly, I aim to empower you with the ability to 
compose live performances of resonant calculation through _your own_ digital instruments - quite literally,
neural synthesizers =)

### 𝑎𝑡𝑙𝑎𝑠 and Recording

Above, I briefly showed you my `std.TemporalBuffer[T]` type which, as it's name suggests, holds a rolling window
of timestamped data.  When a value exceeds the buffer window, it's automatically trimmed out.  By 
default, temporal buffers hold an `atlas.ObservanceWindow` worth of temporal data, but that value can be set with 
a simple `time.Duration` value through configuration.  JanOS can be configured _before_ launch through an extensionless `atlas`
file placed in the same directory as your `go.mod` file, or at runtime through the `atlas` package.  By default,
JanOS will emit a preamble identifying itself on initialization - to disable this, an `atlas` file is required.
The atlas file is currently a standard JSON formatted configuration, but the latest available parameters
are visible by inspecting the `atlas` package in your IDE's typeahead.  The most common configurations to set are _verbosity_
and _silent_ modes, which can be set like so

    // At runtime
    func init() {
        atlas.Verbose(true) // Enables verbose recording
        atlas.Silent(true) // Silences all recordings
	    atlas.ObservanceWindow = time.Second // Sets the default temporal observance window 
    } 

    // In an atlas file
    {
        "verbose": true,
        "silent": true
    }

NOTE: temporal buffers referentially set their observance window, meaning you can change it globally at runtime 

In addition, stdout is considered the output of _recording_ through the `rec` package.  It behaves just as the
`fmt` package, but ensures all recordings include the _name_ of which entity emitted it as a first-order element.  The
concept of _recording_ and _observation_ is heavily used throughout this project - _every_ thought in JanOS is recorded
for others to observe through the neural nexus, which we'll get to later.

### Thoughts

Lastly, I'd like to talk about a neural _thought_

    type Thought[T any] struct {
        Path
        revelation T
        gate       *sync.Mutex
        created    bool
    }
    
    // NewThought creates a new instance of a Thought[T] and assigns it a unique Path identifier by calling id.Next()
    //
    // NOTE: If you'd like to assign the thought's path directly, you may provide it through the variadic.
    func NewThought[T any](revelation T, path ...string) Thought[T] {
        p := []string{strconv.FormatUint(id.Next(), 10)}
        if len(path) > 0 {
            p = path
        }
        return Thought[T]{
            Path:       p,
            revelation: revelation,
            gate:       new(sync.Mutex),
            created:    true,
        }
    }
    
    func (t *Thought[T]) sanityCheck() {
        if !t.created {
            panic("std.Thought[T] must be created through std.NewThought[T]")
        }
        if t.gate == nil {
            t.gate = new(sync.Mutex)
        }
    }

    // Revelation sets and/or gets the Thought's inner revelation.  If no value is
    // provided, it merely gets - otherwise it sets the value before returning it.
    func (t *Thought[T]) Revelation(value ...T) T {
        t.sanityCheck()
        t.gate.Lock()
        defer t.gate.Unlock()
        
        if len(value) > 0 {
            t.revelation = value[0]
        }
        return t.revelation
    }

This structure is _quite_ simple - it simply gates access to the underlying value - but it marks the foundation
of the next major component of JanOS: _the neural nexus!_  A thought is often an _inspirative nudge_ towards 
the neuron's ultimate workload.  While it prompts the neural endpoint with data, the neuron must translate what the thought
means to _itself._  For instance, a neuron could take in _integer_ thoughts which translate into a more complicated
_floating-point_ formula.  The way thoughts are tracked through a neural nexus is by a simple `map[string]any`,
but the data contained within each map entry could be _another_ map.  Each map represents a subgroup of related
data points and should use another string identifier.  This allows us to traverse through the recursive map
structures using a `std.Path` type

    // A Path is a sequence of string identifiers which can act as a chain of custody.
    type Path []string
    
    // Emit outputs the path as a delimited string.
    //
    // NOTE: If no delimiter is provided, ⇝ is used - otherwise, only the first delimiter provided is used.
    func (t Path) Emit(delimiter ...string) string {
        d := "⇝"
        if len(delimiter) > 0 {
            d = delimiter[0]
        }
    
        return strings.Join(t, d)
    }

The path type acts as a recursive map key, where each string identifier should reference another map until
you reach the ultimate value you're looking for.  These are _string keys_ for one important reason: neural
nexuses can cross application boundaries!  JanOS aims to allow _infinitely many_ instances to co-calculate
on the same thoughts in real-time - in fact, that's one of the prime directives of a neural architecture.
In your machinery, each subsystem works as an isolated component in a larger cohesive system -
motor control, perception, thought, proprioception, all the things.  The goal is to allow the very same
mechanic _with ease._

We'll get to the nexus later, for now I feel like I've overloaded you!  I promise this isn't as daunting as it
appears on the surface - as long as you consider _yourself_ a neuron in a larger network, being nudged by the
inspirative thoughts guiding you along _your_ path =)