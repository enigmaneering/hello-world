# `E0S2 - The Path`
### `Alex Petz, Ignite Laboratories, October 2025`

---

### Locating Thoughts
Now that we have ways to _gate_ and _disclose_ a thought, we get to locate it using a path.  A Path is a _very_ simple type

    // A path definition
    
    type Path []any

    func (p Path) String() string { return p.Emit() } 

    func (Path) Emit(delimiter ...string) string {
        d := "⇝"
        if len(delimiter) > 0 {
            d = delimiter[0]
        }
    
        return strings.Join(StringifyMany(p), d) 
    }

_That's it!_  It's just a slice of stringable objects which can be used to reference a location.  If this strikes
imagery of a URI, I will have hit my mark!  The concept is _no different_ - but we aren't just locating _objects_,
we're also locating object _members._  Take, for instance, the below setup

    type AStruct struct {
        Other BStruct
    }

    type BStruct struct {
    }

    func (BStruct) String() string { }

    var MyStruct = AStruct{}

Using a `std.Path` you could describe the `String()` method on the sub-structure as such:

    p := std.Path{"MyStruct", "Other", "String"}

As we'll see in the next step, however, the real power of the path comes from implicitly referencing _map keys._

    type AMap map[string]any

    var MyMap = make(AMap)

    MyMap["KeyA"] = make(map[string]any)
    MyMap["KeyA"]["KeyB"] = "Hello, World!"

    ...
    
    p := std.Path{"MyMap", "KeyA", "KeyB"} // Yields "Hello, World!"

The above creates a relational pathway that steps though `KeyA` on the global `MyMap` variable - which, itself, is a map 
which steps to `KeyB` - describing a path to observing "Hello, World!" at runtime.

That's the point of a path!  Lazy object relational mapping at the _type_ level using the Stringable pattern =)