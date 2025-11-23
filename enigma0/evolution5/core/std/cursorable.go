package std

// A Cursorable entity is one that can traverse an abstract space using relative or absolute motion.
//
// NOTE: Think of the very cursor currently in your IDE being driven between points arbitrarily while selecting data.
type Cursorable[TOut any] interface {
	// Jump performs a relative instantaneous jump 𝑛 positons forwards or backwards and -then- yields the resulting element.
	Jump(n any) Cursorable[TOut]

	// JumpTo performs an absolute instantaneous jump to position 𝑖 and -then- yields the resulting element.
	JumpTo(i any) Cursorable[TOut]

	// JumpAlong will instantaneously move to the result of the "next" function and -then- yields the resulting element.
	// This will continue until 'next' is exhausted or returns nil (if given a function provider.)
	//
	// NOTE: You may alternatively provide a slice of direct TOut points to traverse, or an individual TOut for a single element.
	JumpAlong(steps any, relative bool) Cursorable[TOut]

	// Walk relatively traverses 𝑛 positons forwards or backwards at a rate of 'stride', yielding each element -after- each step.
	Walk(n any, stride any) Cursorable[TOut]

	// WalkTo absolutely traverses to position 𝑖 at a rate of 'stride' and yields each element -after- each step.
	//
	// NOTE: If the final element exists less than the stride distance from the last step, it will still be stepped to and yielded.
	WalkTo(i any, stride any) Cursorable[TOut]

	// WalkAlong will move to the result of the "next" function at a rate of 'stride' and yields each element -after- stepping to it.
	// This will continue until 'next' is exhausted or returns nil (if given a function provider.)
	//
	// NOTE: stride will be revealed between each step, allowing you to make a "dynamic stride" using function providers.
	//
	// NOTE: You may alternatively provide a slice of direct TOut points to traverse, or an individual TOut for a single element.
	WalkAlong(steps any, stride any, relative bool) Cursorable[TOut]

	// Current returns the current position's element.
	//
	// NOTE: This is an entirely -independent- action of the movement operation chain, meaning it should not affect a call to Yield.
	Current() TOut

	// Yield returns the elements found from the current movement operation chain.
	Yield() []TOut
}

/*
# Six Degrees of Semantic Freedom

Cursor operations are meant to be fluently chained.  If you desire only the element where the cursor
currently resides, a call to Current() should do so independently of the fluent chain.  Otherwise,
you would Yield() the results when you're ready to.

Serialization of cursor motivation

Single brackets are used for standard motivation using traditional index access:
[:99] <- Yields 0-99
[42:] <- Yields 42-(len-1)
[42:99] <- Yields 42-99
[42] <- Yields 42

Double brackets indicate cursor motivation.  Moving a cursor requires three aspects:

0 - Is this relative or absolute motion?  (indicated by a '~' prefix of the first operator to indicate 'relative')
1 - Is this a 'jump' or a 'walk'?  (indicated by the presence of a second operator)
2 - If walking, what's the stride? (indicated by the value of the second operator)

NOTE: The prefix changed from '*' to support variables instead of just integers, as a * dereferences inline

For example:

Jumping --
[[~0]] <- Jump(0) - Yields the current element
[[0]] <- JumpTo(0) - Yields the element at position 0

[[~42]] <- Jump(42) - Yields the element 42 positions away
[[~-42]] <- Jump(-42) - Yields the element -42 positions away

[[42]] <- JumpTo(42) - Yields the element 42 positions from 0
[[-42]] <- JumpTo(-42) - Yields the element -42 positions from 0

Walking --
[[~42, 0]] <- Walk(42, 0) - Yields nothing, as a zero stride yields zero elements

[[~0, 1]] <- Walk(0, 1) - Yields the current element
[[0, 1]] <- WalkTo(0, 1) - Yields the elements from here to "position 0" at a stride of 1

[[~42, 1]] <- Walk(42, 1) - Yields the elements from here to "42 positions away" at a stride of 1
[[~-42, 1]] <- Walk(-42, 1) - Yields the elements from here to "-42 positions away" at a stride of 1

[[ 42, 1]] <- WalkTo(42, 1) - Yields the elements from here to "position 42 from zero" at a stride of 1
[[-42, 1]] <- WalkTo(-42, 1) - Yields the elements from here to "position -42 from zero" at a stride of 1

[[~42, 5]] <- Walk(42, 5) - Yields the elements from here to "42 positions away" at a stride of 5
[[~-42, 5]] <- Walk(-42, 5) - Yields the elements from here to "-42 positions away" at a stride of 5

[[ 42, 5]] <- WalkTo(42, 5) - Yields the elements from here to "position 42 from zero" at a stride of 5
[[-42, 5]] <- WalkTo(-42, 5) - Yields the elements from here to "position -42 from zero" at a stride of 5

[[ 42, -5]] <- WalkTo(42, -5) - Yields the elements from here to "position 42 from zero" at a stride of -5 (see "bounded contexts")
[[~42, -5]] <- Walk(42, -5) - Yields the elements from here to "42 positions away" at a stride of -5 (see "bounded contexts")

Bounded Contexts -
Most data is naturally bounded by its own size, so cursor movement naturally provides Python-style 'tail indexing.'
In addition, a negative -stride- indicates to traverse the inverse of the shortest path, causing the cursor to
take the "long way" 'round.  That being said, if your cursor is traversing a mathematically infinite space, a
negative stride is an 'undefined' operation that should cause a panic.

This has an important caveat: saturated contexts!  In a "clamped" environment, there's only -one- direction that
you could logically traverse to reach a target - thus, a panic should occur if given a negative stride while clamped.

In a path, these can be chained:

[42][[~42]] <- Yields position 42 and then the element 42 positions away from it
[42][[~99, 5]] <- Yields position 42 and then strides 5 elements at a time towards element 99 until it reaches it.
[42][[~-99, 5]] <- Yields position 42 and then strides 5 elements at a time towards element -99 until it reaches it.
[42][[99, 5]] <- Yields position 42 and then strides 5 elements at a time towards element 99 until it reaches it.
[42][[-99, 5]] <- Yields position 42 and then strides 5 elements at a time towards element -99 until it reaches it.

NOTE: The 'along' operations are simply a shorthand for calling the relative and absolute functions dynamically,
meaning they serialize whatever they are given by emitting a chain of the above rules.  The 'along' operations take
in either a func() ~T (for JumpAlong), a func() (~T, T) (for WalkAlong), or a serialized set of the instructions
as described above.

# Why an interface!?

Great question, Alex!  Because this allows different kinds of cursors!  For example, a "marquee" cursor could build
a cubic box in a three-dimensional environment when given two coordinates to select information.  Or a "lasso" cursor
could be used to generate manifold shapes that select regions of space.  Lastly (at least for now) a "direct" cursor
would be used to directly select the information.

Types:
- Marquee(points) - Higher dimensions above 3 should select single points, with a cube being the innermost
- Lasso(points)
- Radial(radius)
- Direct(points)
- Convolution(points, kernel)
*/
