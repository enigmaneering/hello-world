# `E0S1 - The Idea Nexus`
### `Alex Petz, Ignite Laboratories, October 2025`

---

### Ideas

Now that we've explored the concept of isolated activations, we get to consider how to _persist_ thoughts between
them!  Currently, we've handled this through a referential _thought_ - but that only works if you actually
_have_ the damned pointer!  The thought provides thread-safe access by exposing synchronization over a mutex,
but to _locate_ a thought requires formulating it into a `std.Idea`.  An idea is very much like a URI in the
world wide web - _intentionally so,_ I might add, as it's a wonderful paradigm already - used to _describe_
the thought's location.  Ideas _are_ thoughts, but enhanced with a few extra features.  For instance, 
an idea provides a primitive level of privacy around the thought - handled through a `std.Secret` - allowing only 
those with a _code_ to access it (but we'll get to that later)

Instance-level ideas are housed within the `idea` package and are _globally_ accessible using two methods

    idea.Get[T](path std.Path, code ...any) T
    idea.Set[T](value T, path std.Path, code ...any)

The `idea` package houses the central `std.Idea` that houses all other ideas - essentially, a 
`map[string]any` where each value could potentially be _another_ map.  A std.Path is a way of expressing
a `[]string` type where each element represents a map key - _that's really it!_

The more astute of you will have already noticed the hints of an _Entity Component System_ pattern!  The _entity_
represents everything in the system - neural activations and thoughts - while the idea nexus represents the shared
_components_ and each synaptic cluster is a _system_.  In fact, _The Treaty of Orlando_ (which formalized ECS) was a
huge source of inspiration for this OS.  Academically discussing the concept of _empathy_ and _sharing_ in data structures,
rather than secrecy and encapsulation, was a breath of philosophical fresh air!  Stein, Lieberman, and Ungar were _far_
ahead of their time in a world rapdily changing.