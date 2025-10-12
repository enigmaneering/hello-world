package std

import (
	"time"

	"git.ignitelabs.net/janos/core/std"
)

/*
- NOTE

The 'core' JanOS module is where the std library actually lives, but these solutions will demonstrate how each
neural component was added to the existing library.  In this case, the std.Entity is a basic type that provides
a globally unique ID and Given to every object that composes of it.

	- Alex
*/

type Impulse struct {
	std.Entity
	Bridge
	Activation
	Synapse

	Timeline *std.TemporalBuffer[Activation]
}

func (imp Impulse) String() string {
	if len(imp.Bridge) == 0 {
		return imp.Name.Name
	}
	return imp.Bridge.String() + " ⇝ " + imp.Name.Name
}

// RefractoryPeriod represents the duration between the last activation's completion and this impulse's firing moment.
func (imp *Impulse) RefractoryPeriod() time.Duration {
	last := imp.Timeline.Latest()
	if len(last) <= 0 {
		return 0
	}
	return imp.Activation.Fired.Sub(*last[0].Element.Completion)
}

// CyclePeriod represents the duration between the last and current activation firing moments.
func (imp *Impulse) CyclePeriod() time.Duration {
	last := imp.Timeline.Latest()
	if len(last) <= 0 {
		return 0
	}
	return imp.Activation.Fired.Sub(*last[0].Element.Fired)
}

// ResponseTime represents the duration between inception and firing of the current impulse.
func (imp *Impulse) ResponseTime() time.Duration {
	return imp.Activation.Fired.Sub(imp.Activation.Inception)
}

// RunTime represents the duration between firing and completion of the last impulse's activation.
func (imp *Impulse) RunTime() time.Duration {
	last := imp.Timeline.Latest()
	if len(last) <= 0 {
		return 0
	}
	return last[0].Element.Completion.Sub(*last[0].Element.Fired)
}

// TotalTime represents the duration between inception and completion of the last impulse's activation.
func (imp *Impulse) TotalTime() time.Duration {
	last := imp.Timeline.Latest()
	if len(last) <= 0 {
		return 0
	}
	return last[0].Element.Completion.Sub(last[0].Element.Inception)
}
