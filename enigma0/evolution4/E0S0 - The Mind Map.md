# `E0S0 - The Mind Map`
### `Alex Petz, Ignite Laboratories, October 2025`

---

A mind map is an _abstract_ concept that exists across _any_ intelligent system, no matter the shape or form.  An
intelligent system that can recognize and respond to a request _participates_ in the mind map but can _implement_
it however it so chooses.  A mind map has two halves - internally, `stdmem`, and externally, a `nexus` - we begin
internally

### Neural LIQs
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
    liq := std.Path{"Field", 42}) // a Go LIQ slice

You use a LIQ like you use a URI - to _locate_ a stored value in `std.Memory` - and if the value cannot be
resolved, a LIQ should fail gracefully.  LIQs are _language implemented_ in the sense that each endpoint
must implement it in their own _language,_ but the implementation is less important than the _pattern._  A
LIQ, when transmitted over the wire, must still be serialized, but at runtime can treat the
individual components as comparable _objects_ that simply _identify_ something.  How the data is serialized
doesn't matter, but each endpoint must accept _any_ format and intelligently parse it on the fly - otherwise,
it should gracefully yield an appropriately descriptive error.

    tl;dr - LIQs abstractly describe a path to an object at runtime 

###  stdin, stdout, stderr, and stdmem


### ideas
From this point onward the concepts are _somewhat_ universal, but each language comes with pros and cons.  I've
adopted Go wholeheartedly, so my solutions will be executed as such.  Go comes with a couple of wonderful quirks,
which I've adopted some non-standard solutions for.  Most importantly, Go does not support generic _methods_ on
structural types; as such, I've chosen to use packages to "namespace" my generic needs.  My work is embodied into
a neural operating system I call _JanOS,_ after the Roman God of time, duality, gates, and beginnings - so any
namespacing as such will be in reference to my design

    std.Memory - The origin 