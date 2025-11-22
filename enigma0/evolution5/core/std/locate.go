package std

func Locate(source any, relative Path) (any, error) {
	// 0. Pop the first step in the path
	// 1 Find the type of the source
	// 1.0 If a map, treat the step as a key
	// 1.1 If a slice or array, treat the step as an index accessor - parseable from a string
	// 1.2 If an object, treat the step as a Stringable field or method accessor
	// 1.3 If a pure function, call it if the step is empty - otherwise, panic (we can't dig deeper into that)
	// 1.4 If a providing function, call it and then treat the result as what this step should be addressing
	// 1.5 If another thought, reveal it using the optional code in the current step and treat the result as what this step should be addressing
	// 1.6 If it's a bidirectional or send-only channel, send the step into it - if receive-only, panic
	// 1.7 Otherwise, panic since we cannot evaluate INTO any other type - only reveal it
	// 2. Recursively call this method on the result with the current step removed in the path

	panic("not implemented")
}
