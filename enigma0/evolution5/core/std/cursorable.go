package std

// A Cursorable entity is one that can traverse an abstract space using relative or absolute motion.
//
// NOTE: Think of the very cursor currently in your IDE being driven between points arbitrarily while selecting data.
type Cursorable[TOut any] interface {
	// Jump performs a relative instantaneous jump 𝑛 positons forwards or backwards and -then- yields the resulting element.
	Jump(n any) Cursorable[TOut]

	// JumpTo performs an absolute instantaneous jump to a specific position and -then- yields the resulting element.
	JumpTo(i any) Cursorable[TOut]

	// JumpAlong will instantaneously move to the result of the "next" function and -then- yields the resulting element.
	// This will continue until 'next' is exhausted or returns nil (if given a function provider.)
	JumpAlong(next any) Cursorable[TOut]

	// Walk relatively traverses 𝑛 positons forwards or backwards at a rate of 'stride', yielding each element -after- each step.
	Walk(n any, stride any) Cursorable[TOut]

	// WalkTo absolutely traverses to a specific position at a rate of 'stride' and yields each element -after- each step.
	//
	// NOTE: If the final element exists less than the stride distance from the last step, it will still be stepped to and yielded.
	WalkTo(i any, stride any) Cursorable[TOut]

	// WalkAlong will move to the result of the "next" function at a rate of 'stride' and yields each element -after- stepping to it.
	// This will continue until 'next' is exhausted or returns nil (if given a function provider.)
	//
	// NOTE: stride will be revealed between each step, if you'd like a "dynamic stride."
	WalkAlong(next any, stride any) Cursorable[TOut]

	// Current returns the current position's element.
	//
	// NOTE: This is an entirely -independent- action of the movement operation chain, meaning it should not affect a call to Yield.
	Current() TOut

	// Yield returns the elements found from the current movement operation chain.
	Yield() []TOut
}

/*
Serialization of cursor motivation

Implicit motivation through standard index access -
[:99] <- Yields 0-99
[42:] <- Yields 42-(len-1)
[42:99] <- Yields 42-99
[42] <- Yields 42

Explicit motion access -
[[*42]] <- Jump(42) - Yields the element 42 positions away
[[42]] <- JumpTo(42) - Yields the element 42 positions from 0

[[*42]] <- Walk(42, 1) - Yields the elements from here to "42 positions away" at a stride of 1
[[42]] <- WalkTo(42, 1) - Yields the elements from here to "position 42 from zero" at a stride of 1
[[*42, 5]] <- Walk(42, 5) - Yields the elements from here to "42 positions away" at a stride of 5
[[42, 5]] <- WalkTo(42, 5) - Yields the elements from here to "position 42 from zero" at a stride of 5

In a path, these can be chained:

[42][[*42]] <- Yields position 42 and then the element 42 positions away from it
[42][[*99, 5]] <- Yields position 42 and then strides 5 elements at a time until it reaches an element 99 away.
[42][[*-99, 5]] <- Yields position 42 and then strides -5 elements at a time until it reaches an element -99 away.
[42][[99, 5]] <- Yields position 42 and then strides 5 elements at a time until it reaches element 99.
[42][[-99, 5]] <- Yields position 42 and then strides -5 elements at a time until it reaches element -99.

NOTE: This means directionality is implicit and should follow the shortest path (respecting saturation)
*/
