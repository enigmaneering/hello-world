package std

// A Range describes the slice accessor pattern[low:High]
//
// NOTE: Use 'nil' for either bound to describe an open-ended slice, or for both to describe the entire slice.
type Range[TOut any] struct {
	Low  TOut
	High TOut
}
