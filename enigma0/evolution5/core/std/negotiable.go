package std

// Negotiable describes a type which can receive a code and perform arbitrary advanced "negotiation" against it.  This
// means it can take in -any- type and introspect it to determine if it should "pass" or "fail" a test.
type Negotiable interface {
	// Negotiate should take in a code and interface with it - if it deems it acceptable to "pass", it should yield true.
	Negotiate(code any) bool
}

// negotiable tests if a value is a Negotiable type.
func negotiable(values ...any) bool {
	if len(values) == 0 {
		return false
	}
	allGood := true
	for _, value := range values {
		if _, ok := value.(Negotiable); !ok {
			allGood = false
		}
	}
	return allGood
}
