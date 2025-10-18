package std

import "git.ignitelabs.net/janos/core/std"

type Nexus struct {
	std.Entity
	Idea[map[string]any]
}
