package std

import "time"

// A Heart is used to signal presence and activity.  By itself, it doesn't mean much - but when paired with a
// context, this holds the last moment it was pulsed.  This allows you to periodically signal back to external
// systems that you are still operating as expected when running a long process.
type Heart time.Time

// Beat sets the last pulse moment of the Heart to time.Now.
func (h *Heart) Beat() {
	*h = Heart(time.Now())
}
