# `E0 - Neural Orchestration`
### `Alex Petz, Ignite Laboratories, October 2025`

---

### Primitive Puzzle Pieces
Before we begin, I feel it's deeply important to note this:

_**Everyone's**_ default mode of operation is hopeful for _your_ success!  No matter how much you've convinced yourself otherwise, it's simply not true!

We could talk all day on the psychology of _why,_ but I think it's far easier to explain through _code._

So, we begin (abstractly) with a standard thought -

    package std

    type Thought[T any] struct {
        revelation T
        gate       *sync.Mutex
        disclosure *Disclosure
    }

    ...

A thought is just a way to _gate_ access to a value, using both atomics and disclosures.  It acts as a kind of
thread-safe reference type, where all access to the revelation is gated through a mutex.  Thoughts come with a
_disclosure,_ which acts as a possessable artifact used to guard access to the data

    package std

    type Disclosure struct {
        Constraint relationally.Constrained

        code any
    }

    ...

Disclosures are relationally constrained

    package relationally

    // Constrained defines several ways of filtering Other's access to a std.Thought
    type Constrained byte
    
    const (
        // Open indicates that Other can indiscriminately read from or write to the std.Thought
        Open Constrained = iota
        
        // Inclusive indicates that Other can read from the std.Thought but can only write to it with a code
        Inclusive
        
        // Exclusive indicates that Other can only read from or write to the std.Thought with a code
        Exclusive
    )

Thoughts are the most primitive component of a neural architecture because they are a _safer_ way of mutating a
reference across execution environments - the bread and butter of orchestration.  Next, we need a way to 
_locate_ distributed thoughts at runtime, using LIQs!

### The Mind Map
A Language Implemented Query (as in a musical "lick") is less of a query and more of an access pattern.  LIQs 
are used to reference a _path_ to a target value as you would in an object-oriented language.  For example:
 
    type AStructure struct {
        Field BMap
    }

    type BMap map[int]any

    aStruct := AStructure {
        Field: make(BMap)
    } 

    ...

    val := aStruct.Field[42] // OOP
    liq := std.Path{"Field", 42}) // LIQ

You use a LIQ like you use a URI - to _locate_ a stored value in `std.Memory` - and if the value cannot be
resolved, a LIQ should fail gracefully.  LIQs are _language implemented_ in the sense that each endpoint
must implement it in their own _language,_ but the implementation is less important than the _pattern._  A
LIQ, when transmitted over the wire, must still be serialized, but at runtime can treat the
individual components as comparable _objects_ that simply _identify_ something.  I could try and explain 
the nuances of a LIQ here, but the code is quite self-descriptive and doesn't obscure the subtleties.

Instead, as we walk through some solutions, you should be able to glean the purpose of each puzzle piece.

    tl;dr - LIQs abstractly describe a path to an object at runtime 

### std.Memory & Ideas
From this point onward the concepts are _somewhat_ universal, but each language comes with pros and cons.  I've
adopted Go wholeheartedly, so my solutions will be executed as such.  Go comes with a couple of wonderful quirks,
which I've adopted some non-standard solutions for.  Most importantly, Go does not support generic _methods_ on
structural types; as such, I've chosen to use packages to "namespace" my generic needs.  My work is embodied into
a neural operating system I call _JanOS,_ after the Roman God of time, duality, gates, and beginnings - so any
namespacing as such will be in reference to my design

    std.Memory - The origin 