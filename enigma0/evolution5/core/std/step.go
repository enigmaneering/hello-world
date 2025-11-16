package std

// A Step is a single point in a Path that describes an endpoint and an optionally necessary code for accessing it.
type Step struct {
	Data any
	Code any
}
