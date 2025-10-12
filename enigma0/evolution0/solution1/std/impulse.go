package std

import (
	"time"

	"git.ignitelabs.net/janos/core/std"
	"git.ignitelabs.net/janos/core/sys/given/format"
)

type Impulse struct {
	std.Entity
	Bridge
	Activation
	Synapse

	Beat uint

	Timeline *std.TemporalBuffer[Activation]

	Thought *Thought
}

func NewImpulse(named ...string) *Impulse {
	
	imp := &Impulse{
		Entity: std.NewEntity[format.Default](),
	}

}

func (imp Impulse) String() string {
	if len(imp.Bridge) == 0 {
		return imp.Name.Name
	}
	return imp.Bridge.String() + " ⇝ " + imp.Name.Name
}

// BridgeTo checks if the target impulse's bridge seeds from the source impulse's
func (imp *Impulse) BridgeTo(target *Impulse) {
	source := imp.Bridge
	if source == nil || len(source) == 0 {
	}
	source = append(source, imp.Name.Name)

	alreadySet := true
	for i, b := range source {
		if len(target.Bridge) < len(source) || i >= len(target.Bridge) {
			alreadySet = false
			break
		}
		if target.Bridge[i] != b {
			alreadySet = false
			break
		}
	}
	if !alreadySet {
		target.Bridge = append(target.Bridge, source...)
	}
}

// Manifest either returns the current Thought (if not nil), or sets it to the provided 'default' value and returns it.
func (imp *Impulse) Manifest(nilValue any) *Thought {
	if imp.Thought == nil {
		imp.Thought = NewThought(nilValue)
	}
	return imp.Thought
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
