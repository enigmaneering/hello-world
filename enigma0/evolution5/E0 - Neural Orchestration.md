# `E0 - Neural Orchestration`
### `Alex Petz, Ignite Laboratories, October 2025`

---

We begin by asking a question.

### Who are you?

From this thread's perspective?  You're _Other!_

The Yin to the Yang, the Spock to Captain Kirk, and the rhythm to the blues.

Without you or I, this groove wouldn't jive and the axis of Life wouldn't twist and twirl.

So, I'm really excited to tell you this:

_**Everyone's**_ default mode of operation is hopeful for _**Others'**_ success!  

No matter how much you've convinced yourself otherwise, it's simply not true!

We could talk all day on the psychology of _why,_ but I find it's far easier to demonstrate the perpetual
motion of empathy through _**code**_ =)

### Primitive Puzzle Pieces

Before anything can happen, a standard thought must form -

    package std

    type Thought[T any] struct {
        revelation T
        gate       *Gate 
        disclosure *Disclosure
    }

    ...

Think of a thought as a way to _give_ a value and then _disclose_ who can access it while _gating_ the inquiring
entities into a single file line.  A _gate_ is a way of **patiently** attempting to access a shared point, while a 
_disclosure_ describes the secret code which **regulates** access to the thought.

    package std

    // A Gate can patiently Attempt to interface with a sync.Mutex.
    type Gate struct {
        sync.Mutex
    }
    
    // Attempt repeatedly tries to "hold the gate" before conceding if unable to do so in the allotted time.
    //
    // NOTE: This will return true if a lock was attained - otherwise, false.
    func (g *Gate) Attempt(timeout time.Duration) bool {
        ...
    }

    // A Disclosure describes the codified conditions and constraints behind Others' access to a Thought.
    type Disclosure struct {
        Constraint relationally.Constrained

        code any
    }

    ...

A disclosure allows compartmentalization of who can access the thought with a relational _constraint_

    package relationally

    // Constrained defines several ways of filtering Others' access to a std.Thought
    type Constrained byte
    
    const (
        // Open indicates that Other can indiscriminately read from or write to the std.Thought
        Open Constrained = iota
        
        // Inclusive indicates that Other can read from the std.Thought but can only write to it with a code
        Inclusive
        
        // Exclusive indicates that Other can only read from or write to the std.Thought with a code
        Exclusive
    )

**KEY: By default, all thoughts are universally "open" to Other, making constraint a _choice_.**  This reflects 
Edmondson’s view that openness is the foundation of psychological safety, and it bolsters Erikson’s concept of 
generativity as stagnation becomes a deliberate _choice_ in a "live" architecture.  More consequentially, this design 
facilitates Jungian synchronicity - by disconnecting the storage and recollection of thoughts into a 
distributed substrate, archetypal patterns can surface collectively, and disparate systems can 
intelligently nudge the same shared thought into meaningful alignment.

    tl;dr - open thoughts are what facilitate synchronicity

**Introspection** is the most primitive _aspect_ of a neural architecture because it defines a _safer_ way of mutating a
reference across execution environments - the bread and butter of neural orchestration.  This parallels Life, as reflection
often reveals the true nature of what's _actually_ affecting our existences - both positively and negatively.  

Next, we get to define a way to _locate_ distributed thoughts at runtime: using LIQs!

### Mental Licks
A Language Implemented Query (or LIQ, as in a _musical_ "lick") is less of a query and more of an access pattern.  Licks 
are used to reference a _path_ to a target value as you would in an object-oriented language.  Think of them as a
kind of "mind map" describing exactly how to traverse _anything_ to reach the treasure.  Rather than redefining the very
same access patterns programmers already love and adore, a lick _**literally**_ redescribes it!

For example:
 
    type AStructure struct {
        Field BMap
    }

    type BMap map[CType]any

    aStruct := AStructure {
        Field: make(BMap)
    } 

    ...

    val := aStruct.Field[42] // OOP from 'aStruct' -specifically-

    liq := std.Path{"Field", 42} // A LIQ relative from -anything- (in JanOS, a std.Path)
    val := mem.Recall[CType](aStruct, liq...) // Revealing the lick from 'aStruct' specifically

NOTE: All memory is stored in `std.Memory` - but, as Go does not support generic _methods,_ typed access is available through
_functions_ in the `mem` package.

You use a lick as you would a URI to relatively _locate_ a value in `std.Memory` - if the value cannot be
resolved, a lick should fail gracefully.  Licks are _language implemented_ in the sense that each endpoint
must implement it in their own _language,_ but the implementation is less important than the _pattern._  A
lick, when transmitted over the wire, must still be serialized, but instead of serializing _data_ you're 
serializing a _path_ to something.  For instance, if multiple systems are simulating a model of physical
space, a `std.Path{42,"π",11.22}` could universally be understood as a _typeless_ `(x,y,z)` coordinate. This 
allows you to reference space _**as a concept**_ - rather than needing to know _anything_ about the
system's _**specific implementation**_ of a spatial simulation.

I could try and explain the nuances of a lick here, but the code is quite self-descriptive and doesn't 
obscure the subtleties. Instead, as we walk through some solutions, you should be able to glean the purpose 
of each lick.

    tl;dr - LIQs abstractly describe a path to anything at runtime

From this point onward the concepts are _somewhat_ universal, but each language comes with pros and cons.  I've
adopted Go wholeheartedly, so my solutions are almost exclusively written in it.  My work is embodied into a 
neural operating system I call _JanOS,_ after the Roman God of time, duality, gates, and beginnings. 

### Empathy and the Mind Map
I'd like you to take a moment to read a _pivotal_ piece of literature I've included here, called `A Shared View of 
Sharing: The Treaty of Orlando`.  In it, Stein, Lieberman, and Ungar (1988) articulate a foundational issue with 
object-oriented design: the ability to empathetically _evolve_ a programmatic structure (p7-8).  While their ideas, at the 
time, were confined to _language_ design - their concepts hold truer when focusing on _system_ design.  The ability 
for two systems to understand how to recall an _aspect_ of something, rather than knowing _every detail_ of it, is crucial - and 
that's _exactly_ what the mind map aims to solve.

    tl;dr - how can a system empathetically share itself without remaining vulnerable?

It'd be quite naïve to expose every aspect of a system unabashedly and without disclosure, hence the reasoning for such
a mechanic as described above.  So, what's _different_ about working with a mind map?  Let's take a look at how to
access memories:

    // Remember: the `mem` package provides typed access to `std.Memory`

    data := std.Memory.Reveal("A", "Path", 2, "Something") // Reveals "A.Path[2].Something" as an 'any' type
    data :=  mem.Recall[TOut]("A", "Path", 2, "Something") // Reveals "A.Path[2].Something" as 'TOut'

Your structure (or "entity") should always be available by logically walking a path.  This, in turn, allows
many unique structures (or "entities") to implicitly _follow_ the same path.  A basic example of this is 
through the `tiny.Operation` type, which uses this mechanic to reference an implicit arithmetic "context" 
which defines the current base and precision.  This is its _entire_ structural definition:

    // An Operation is the source of a calculation chain, beginning with a path to a tiny.Context.
    // Most operations yield a Formula, which can be used to evaluate the operation into a result.
    //
    // see.BaselessCalculation
    type Operation[T any] std.Path

    // NOTE: tiny.Context defines the mathematical radix and precision to use in the operational chain (among other things)

The magic comes in the fluent method chain it seeds from the described path.  When you call a 𝑡𝑖𝑛𝑦 operation, it
generates an anonymous function which operates on whatever context it finds when invoked.  This allows
ad-hoc code to orient itself at execution with minimal effort on the engineer's part.

What do I mean by "orient itself"?  Well, let's add the numbers 46 and 2:

                         ⬐ Creates a base₁₀ (and default precision) tiny.Operation[int]
    result, err := tiny.Int.Add(46, 2).Equals() ← Lazily yields an (int, error)
                             ⬑ Creates a tiny.Formula[int]

Mathematically, the underlying _type_ of the operands doesn't really matter (as long as it can be converted into
digits) - but, the _**act of describing a path**_ to a standard `tiny.Context` allows _this specific operation_ to
"inject" _**any**_ type emergently.  For instance, if you'd like to track modular arithmetic features (like the modulus,
remainder, and breach), `TOut` can embody `tiny.Modular` - causing your operational chain to speak _in that 
language._ If you'd like to store and recall _your own_ `tiny.Context`, you can immediately begin a fluent calculation 
chain from the _path_ to wherever you stored it:

    path := std.Path{"Some, Location"}

    resultF, err :=           tiny.Operation[float64](path).Add(46, 2).Equals()
    resultI, err :=               tiny.Operation[int](path).Add(46, 2).Equals()
    resultC, err :=        tiny.Operation[complex128](path).Add(46, 2).Equals()
    resultM, err := tiny.Operation[tiny.Modular[int]](path).Add(46, 2).Equals()

This can be taken one step further and given a factory function to generate your path:

    func CustomOperation[TOut any]() tiny.Operation[TOut] { 
        return tiny.Operation[TOut](std.Path{"Some", "Location"})
    }

    resultF, err :=           CustomOperation[float64].Add(46, 2).Equals()
    resultI, err :=               CustomOperation[int].Add(46, 2).Equals()
    resultC, err :=        CustomOperation[complex128].Add(46, 2).Equals()
    resultM, err := CustomOperation[tiny.Modular[int]].Add(46, 2).Equals()

### Recap

So, to recap the above example, we're empathetically sharing the mathematical _concept_ of arithmetic bounded to contextual
rules while not enforcing strict types upon the engineer.  While 𝑡𝑖𝑛𝑦 implicitly expects `tiny.Numeric` types, your system
could expect _anything._  Deserialization comes to mind as an excellent use case, for example:

    " A Contrived Deserializer API "

    import "custom"

    for core.Alive() {
        msg, err := custom.Deserializer[MyObj](std.Path{"A","Channel"}).Reveal()
        ...
        // handle 'msg'
    }

In this case, your path could literally describe a _channel_ which `custom.Deserializer` blocks for data on before
attempting to serialize it into `MyObj` when a message arrives - creating a neurally impulsed _signal_.  Multiple
subscribers could be used in a kind of "thread pool" to handle messages from the same point concurrently and with
minimal effort.

By centralizing the _data,_ but writing code with its own _perspective_ of it, abstract _concepts_ can be shared
and applied to a multitude of situations =)