# `E0S1 - An Idea`
### `Alex Petz, Ignite Laboratories, October 2025`

---

### Disclosed Thoughts
For neural architectures to work, you'll have to let go of one of your most prized assets: the _function parameter._
While it's a favorite of mine, in neural activation _everything is anonymous._  We still get to leverage _generics_ 
and gain compile time checks of our code - but _every program is unique!_  By extension, every _neuron_ is also 
unique!  Hell, every _function_ is unique!

So, how could you invoke random functions on _**impulse**_ if their signatures aren't standardized!?  

Well, the _engineer_ already knows how to call the target but the _compiler_ needs guidance in _**locating**_ it.  At 
this point, we have thread-safe access - so, let's evolve the fleeting thought type into a something locatable 

    // An idea definition

    package std

    type Idea[T any] struct {
        path        Path
        thought     Thought[T]
        disclosure *Disclosure
    }

    func NewIdea[T any](Thought[T], Path, *Disclosure) (Idea[T], *Disclosure) { }

    // Path simply implements the std.Pathable interface
    func (id Idea[T]) Path() Path { return id.path }

    // Reveal attempts to return the thought's revelation
    func (Idea[T]) Reveal(...code) (T, error) { }

    // Describe attempts to set the thought's revelation
    func (Idea[T]) Describe(T, ...code) error { }

    // Relative attempts to locate an anonymous relative target from this idea
    func (Idea[T]) Relative(path, ...code) (any, error) { }

An idea _is_ a thought, just paired with a _location_ in the form of a `Path` - but we'll get to that in the next
step.  For now, let's talk about _disclosures._  When sharing an idea with another individual, you enter into an
_implicit_ contract of the idea remaining relationally _open._  If the content should remain confidential, the
idea would be considered either relationally _inclusive_ or _exclusive._  A disclosure implements this mechanic

    // A disclosure definition

    package std

    type Disclosure struct {
        Constraint relationally.Constrained
        code       any
    }

    // Relationality

    package relationally

    type Constrained byte
    const (
        Open Constrained = iota
        Inclusive
        Exclusive
    )

When you _create_ an idea, you can either provide it a disclosure or let it create one for you.  The disclosure
can then be used to _gate_ access to the idea's revelation behind a shared "code" using the below rules:

0. `relationally.Open` - The idea is freely accessible for reading and writing without a code 
1. `relationally.Inclusive` - The idea is freely accessible for reading, but requires a code for writing
2. `relationally.Exclusive` - All reading and writing of the idea requires a code

The code is compared at runtime through the _Stringable_ system by calling `Stringify` on each comparator. 

### Stringable and Stringify
Most types can be implicitly stringed, even without the presence of a `String() string` method.  (Numbers especially
so, though that's a _**very**_ deep topic we'll discuss in a future enigma about the `𝑡𝑖𝑛𝑦` module)

Because of this, JanOS offers several ways to parse anonymous data into strings - notably, four important methods:

    // Stringable safely checks if JanOS -can- stringify the parameters
    std.Stringable(any) bool
    std.StringableMany(...any) bool
    
    // Stringify attempts to implicitly string the parameters
    std.Stringify(any) string
    std.StringifyMany(...any) []string

These four methods will implicitly convert your type into a string using a standardized set of rules defined on their
documentation.  For the latest details of how it works, please refer to the code's documentation directly.