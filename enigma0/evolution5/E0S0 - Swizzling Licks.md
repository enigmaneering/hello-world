# `E0S0 - Swizzling Licks`
### `Alex Petz, Ignite Laboratories, November 2025`

---

I've got another question for you

### What are you looking for?

Most of the time, you can glean what you need from the data - but you might need it in a different logical _order,
type,_ or even to be _translated!_

Luckily for us:

**This isn't a new issue!**

In fact, it's so deeply rooted in software development that Darwin (the core of macOS) embraced the concept of swizzling
from Objective-C during the late 1980s - and still uses it to this day!  The term, however, was coined by the graphical
design community as a shorthand for quickly referencing a set of vector values as such:

    g, r, b, a = someVal.GRBA // Abstractly "swizzling" out green, red, blue, and THEN alpha

What I'm trying to say is that, philosophically speaking, Humanity _in general_ has already brooded on this idea and
dreamt for a better implementation of it since the dawn of computing.  The _Objective-C_ community wanted the ability to
dynamically replace _methods_ at runtime - more colloquially, "monkey patching" - but the _graphical_ community wanted 
the ability to merely mix up vectors in a minimal amount of text.  On the surface, they are _no different_ from one
another - how do you dynamically replace a member at runtime, if you're aware of what it means to do so?

A `std.Path` - from another perspective - _is_ a vector of _mixed_ types.   In Go, those types can _include_ method
calls - meaning you could replace method calls _dynamically_ and without issue.  This makes it a prime candidate for 
swizzling =)

But what does it even mean to swizzle out **_anything!?_**

Well, let's say you'd quickly like to build a tuple of a window's coordinates and whether the window has focus - traditionally, this would be your code:

    type CustomTuple struct {
        X uint
        Y uint
        Focused bool
    }

    func (w *Window) Detail() Tuple {
        ...
        return Tuple {
            X: w.X,
            Y: w.Y,
            Focused: w.HasFocus(),
        }
    }

While that may appear _benign,_ let's consider the implications:

    0. A new type named "CustomTuple" has been introduced
    1. You still have to explicitly call each component to address it on a separate line
    2. The "Detail()" method must exist to access the Tuple type
    3. Every ad-hoc "tuple" you make needs a self-descriptive name and must be maintained

_But there's a better way!_

    var w *Window
    ...
    []any {x, y, focused} = swizzle(w, X, Y, HasFocus())

I'd like to introduce you to a new Go keyword - `swizzle` - which takes in several parameters.  The first is the
target to perform swizzling off of, and the remaining are a variadic of named identifiers.  If the identifier is
simply _named,_ it's returned as is - including methods!  Meaning, if you'd like to _not_ invoke `HasFocus()` and
return a callable _function_ to it - you can do exactly that

    []any{x, y, focusFn} = swizzle(w, X, Y, HasFocus)

If we change the identifier from invocation to identification, the third returned parameter is an _anonymous function_
sharing the same signature and pointing to the target.  Under the hood, `swizzle` simply acts as shorthand to produce
the boilerplate code it defines:

    // Your code:
    []any{x, y, focused} = swizzle(w, X, Y, HasFocus())

    // The psuedocode code Go compiles behind the scenes:
    x := w.X
    y := w.Y
    focused := w.HasFocus()

    ---------------------------
    // The function example:
    []any{x, y, focusedFn} = swizzle(w, X, Y, HasFocus)

    // Expanded:
    x := w.X
    y := w.Y
    focusedFn := w.HasFocus

Most importantly, however, you can still use _index accessors_ in the swizzling pattern!

    var data []any
    []any{first, second, third} = swizzle(data, [0], [1], [2])

This is critical because a `std.Path` literally _is_ a `[]any` - meaning you could replicate the concept of "monkey
patching" by essentially "remixing" the steps of a path.  

In that same vein, the parameters _aren't_ required to be members of the **target!**

    var pixel Pixel
    var w *Window
                                                      ⬐ But you can still reference any variable 
    []any{r, g, b, focused} = swizzle(pixel, R, G, B, w.HasFocus())
                                        ⬑ The first parameter is the "target" to swizzle

If the parameter shares the name with a member of the target, the _target's member_ is **always** selected.  In those
circumstances, you could create a local variable to rename access to the desired non-target member.

This is a _very powerful pattern,_ when you combine it with the next topic -

### Six Degrees of Semantic Freedom

JanOS also implements one other majorly important feature: the _cursor!_  A cursor is anything that implements
the `std.Cursorable` interface - which literally defines _six_ degrees of semantic freedom, broken into two
distinct operations with three modes of traversal each:

                       Method |  Motion  | Description
                      Jump(𝑛) | Relative | Instantly jump forwards/backwards 𝑛 positions
                    JumpTo(i) | Absolute | Instantly jump to position 𝑖
       JumpAlong(steps, bool) |          | Instantly jump along the steps absolutely (false) or relatively (true)
                   Walk(𝑛, s) | Relative | Walks 𝑛 positions forwards/backwards at a stride of 𝑠
                 WalkTo(i, s) | Absolute | Walks to position 𝑖 at a stride of 𝑠
    WalkAlong(steps, s, bool) |          | Walks along the steps absolutely (false) or relatively (true) at a stride of 𝑠

_Please_ digest the above!  This defines the general _concept_ of fluent motion!  Every single table outlined below
is a derivative of the above _six_ operations.

A cursor, when implementing the `std.Cursor` interface, should _lazily enqueue_ the motion commands until a call to 
`Yield()` - or, non-trackingly (a "peek") return the current element it's located at using `Current()`.  The steps defined 
above, however, can be expressed in a much friendlier way - similar to index accessors, but in what I call the _"cursor 
accessor"_ pattern.  Cursor accessors provide _compile time_ instructions for **intelligently** querying data.

The cursor accessor is **entirely** _syntactic sugar_ in implementation - **by design!**  There's no reason to
re-invent the way we _loop_ through a data set - instead, we get to define how to logically _traverse_ it!

    tl;dr - a cursor logically selects data

Think of it like a cursor in your IDE jumping around to 'select' blocks of code while processing it

Cursor accessors and swizzling are baked into the version of Go which powers JanOS, but I eventually hope the 
language might embrace this design and include it in for everyone to play with by default!

Let's take a look at a few cursor accessors:

    data[[42]]    ← Jump(42)
    data[42]      ← JumpTo(42)
    data[[42, 4]] ← Walk(42, 1)
    data[42, 4]   ← WalkTo(42, 1)

Cursor access _evolves_ the index accessor pattern without breaking its existing functionality.  Before we 
get into the details of how, you'll notice we're missing two of our degrees of freedom!  The _Along_ operations 
are implemented by _chaining_ the operations together:

    data[[42, 1]][[11]] ← Walk(42, 1) then Jump(11)

Since that's a lot of brackets, you may break up each operation with a `-` character for readability:

    data[[42, 1]]-[[11]]-[22] ← Walk(42, 1) then Jump(11) then JumpTo(22)

The next part of the cursor accessor is that it provides _three_ kinds of brackets -

**[ Square Brackets ]** - Panic when accessing outside the data's boundaries (traditional Go functionality)

**| Pipe Brackets |** - Clamp the movement to the nearest boundary gracefully

**< Angle Brackets >** - Over or underflow the movement to the other side of the data if out of bounds ("flowed" movement)

When the brackets are _doubled-up_ that indicates every movement is _relative_ to the resulting position
of the last operation (starting from an implicit '0' position) - facilitating _**fluent**_ motion.
    
    Operation |   Mode   |  Out of Bounds  | Method
     [42]     | Absolute |      Panic      | JumpTo(𝑖)
     |42|     | Absolute |      Clamp      | JumpTo(𝑖)
     <42>     | Absolute | Over/Under Flow | JumpTo(𝑖)
     [42, 4]  | Absolute |      Panic      | WalkTo(𝑖, 𝑠)
     |42, 4|  | Absolute |      Clamp      | WalkTo(𝑖, 𝑠)
     <42, 4>  | Absolute | Over/Under Flow | WalkTo(𝑖, 𝑠)
    [[42]]    | Relative |      Panic      | Jump(𝑛)
    ||42||    | Relative |      Clamp      | Jump(𝑛)
    <<42>>    | Relative | Over/Under Flow | Jump(𝑛)
    [[42, 4]] | Relative |      Panic      | Walk(𝑖, 𝑠)
    ||42, 4|| | Relative |      Clamp      | Walk(𝑖, 𝑠)
    <<42, 4>> | Relative | Over/Under Flow | Walk(𝑖, 𝑠)

A wonderful feature of this is that _"flowed"_ access is similar to _tail indexing_ in Python!  If you'd like to grab
the last element, a call to `data<-1>` tells the system to yield the element -1 positions from 0 while
observing underflowing.  At compile time, this evaluates to the _intention:_

    out := data[len(data)-1]

In addition, cursor accessors fully support _ranged_ access:

    // NOTE: Doubled up brackets are -fluently- relative!
    
      Operation  |   Mode   |  Out of Bounds  | Method
     [42:99]     | Absolute |      Panic      | JumpTo(𝑖) - Yield data[42:99] as expected with potential panics
     |42:99|     | Absolute |      Clamp      | JumpTo(𝑖) - Same as above, but stop at the data's boundaries and return that element
     <42:99>     | Absolute | Over/Under Flow | JumpTo(𝑖) - Same as above, but overflow or underflow at the data's boundaries
     [42:99, 4]  | Absolute |      Panic      | WalkTo(𝑖, 𝑠) - Similarly, just at a stride of 4
     |42:99, 4|  | Absolute |      Clamp      | WalkTo(𝑖, 𝑠) - Similarly, just at a stride of 4
     <42:99, 4>  | Absolute | Over/Under Flow | WalkTo(𝑖, 𝑠) - Similarly, just at a stride of 4
    [[42:99]]    | Relative |      Panic      | Jump(𝑛) - Inclusively yield the relative elements starting from 42 away and then ending 99 away from that (panicking if out of bounds) 
    ||42:99||    | Relative |      Clamp      | Jump(𝑛) - Same as above, but stop when you reach the data's boundaries and return that element
    <<42:99>>    | Relative | Over/Under Flow | Jump(𝑛) - Same as above, but overflow or underflow at the data's boundaries
    [[42:99, 4]] | Relative |      Panic      | Walk(𝑖, 𝑠) - Similarly, just at a stride of 4
    ||42:99, 4|| | Relative |      Clamp      | Walk(𝑖, 𝑠) - Similarly, just at a stride of 4
    <<42:99, 4>> | Relative | Over/Under Flow | Walk(𝑖, 𝑠) - Similarly, just at a stride of 4

If you look closely at the doubled-up bracket operations, fluent relative motion can be, well, clunky - at 
first!  Luckily, that's why you can _mix_ the brackets if you'd still like to reference an absolute position 
instead.

    // Relatively walk 42 elements forward, then inclusively yield the elements to absolute element data[99] at a stride of 4
    [[42:99, 4]

Yes, this also means you can _mix the boundary functions!_

    // Relatively walk 42 elements forward while over/underflowing the boundaries, then yield the elements 
    // from there to data[99] while panicking if crossing the data's boundary - all at a stride of 4
    <<42:99, 4]
    
But you also have one additional feature - instead of a `:` you can use an `=` to indicate traversal of the range
_"the long way 'round"_ - which is _exclusive_ of the provided points.  Think of it as a way of saying "please select 
everything BUT the elements in this range" - a kind of "full outer join," in our SQL friends' terminology!

     Operation |   Mode   |  Out of Bounds  | Method
      [42=99]  | Absolute |      Panic      | JumpTo(𝑖) - PANIC!!! Infinity is undefined =)
      |42=99|  | Absolute |      Clamp      | JumpTo(𝑖) - From 42, tries to reach 99 through decrementing, but only yields [41, start] because it's clamped
      <42=99>  | Absolute | Over/Under Flow | JumpTo(𝑖) - Same as above, but underflows and continues to include [end, 100]
     [[42=99]] | Relative |      Panic      | JumpTo(𝑖) - PANIC!!! Infinity is undefined =)
     ||42=99|| | Relative |      Clamp      | JumpTo(𝑖) - From the element 42 away, tries to reach an element another 99 away through decrementing, but clamps at 0 - yielding [41, start]
     <<42=99>  | Relative | Over/Under Flow | JumpTo(𝑖) - Same as above, but underflows and includes [end, 100] (NOTE: The single ending bracket!)
     <<42=99>> | Relative | Over/Under Flow | JumpTo(𝑖) - Same as above, but includes the region above the element 99 away from the element 42 away, instead (remember: fluent motion!)

**Order Matters!**

Cursoring implies _directionality_ of travel - especially with _fluent_ motion!  That means that `[44:22]` will
yield the elements _starting_ at 42 and _stopping_ at 22 _**in descending order**._

Because this is just a kind of 'syntactic sugar', you can also use _**variables**_ as your operands!  That means you can define 
the relative motion of how to _dynamically_ "select" a region of space - without reflection and entirely through compiled 
code.  Earlier I emphasized this is a form of _**intelligent**_ querying, and that's exactly what I meant:

    fn := func() (i int) { 
        ... // Calculate an index
        return i
    }
    v := data[42:fn, 4] // Walk from element 42 to the -runtime- result of 'fn()' at a stride of 4

But we haven't even gotten to the _**coolest**_ part, yet - _predicates!_  These are functions that get called with
each element _**during traversal**_ to either mutate or filter it.  At any point in a chain, you can insert a trinity
of functions for how to process the _currently enqueued_ movements.

    // Cursor access with predicates
    [22]-[42:99, 4]-[SelectFn, ForEachFn, TransformFn]

    // A SelectFn Signature - returns true to 'include' the element
    func(T) bool 

    // A ForEachFn Signature
    func(T)
    
    // A TransformFn Signature - returns a mutated version of the element
    func[TOut any](T) TOut

Any of the functions can be provided `nil` to be ignored, or omitted entirely.  These functions are called
_during traversal,_ meaning the aspect of _distinctness_ has to be addressed.  During a single movement,
it's not possible to visit an element twice - but _chained_ movements absolutely can.  By default, the predicate
functions are called _on every single visitation of an element._  With selection, this means you could produce
a data set with duplicate elements - _by design!_  If, instead, you'd like your predicates to be visited
'distinctly' (in SQL terminology), simply prefix those you wish to do so with a `!` character.

    // Cursor access with distinct predicates
    [22]-[42:99, 4]-[!SelectFn, ForEachFn, !TransformFn]

If a _TransformFn_ is present, the output type of the chain _from that point on_ changes from `T` to `TOut` (as 
defined on the provided function signature).

All predicate function _**parameters**_ can be presented in one of four ways.    

    // A ForEachFn with an anonymously typed function signature
    printer := func(element any) {
        fmt.PrintLn(element)
    }

    // A ForEachFn with a typed function signature matching the data set
    printer := func(element MyType) {
        fmt.PrintLn(element)
    }

    // ForEachFns with a variadic of either
    printer := func(element ...any) {
        fmt.PrintLn(element[0])
    }
    printer := func(element ...MyType) {
        fmt.PrintLn(element[0])
    }

While I show _ForEachFns_ in the above example, the same _**parameter**_ aspect applies to _any_ of the predicate
function signatures.

Predicates, in the sense of a `std.Cursor`, represent a call to `Yield()` - providing a new data set for the next
movement operations to be chained off of.  You can chain as many of these operations together as you'd like:

    [42]-[SelectFn]-<<:99, 4>>-|11|-[nil, nil, Transform[int]] // Ultimately yields an 'int' type

Finally, a special operation is reserved for _emitting_ the LIQ as a `std.Lick`:

    var data []any
    l := data<<42=99, 4>>-[!SelectFn]-[emit]

The act of "emitting" the lick diverges the defined path from its original target, while allowing you
to still `Yield()` from it if you hold the original `std.Lick` object.  To create a `std.Lick` you `std.Parse(anything)` -
which can also intelligently _deserialize_ a lick from many input formats.  To _serialize_ a lick, you `Produce() any` it
into an object structure which _can_ be serialized.