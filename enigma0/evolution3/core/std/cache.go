package std

import (
	"git.ignitelabs.net/janos/core/std"
	"git.ignitelabs.net/janos/core/sys/given/format"
)

// Memory is the standard memory cache embedded into every JanOS instance.
//
// NOTE: For typed memories, please see the 'memory' package directly as Go currently doesn't support generic methods =)
var Memory = NewCache(Path{}, "origin")

// A Cache represents a named Idea map.
type Cache struct {
	Entity std.Entity

	memories Idea[map[string]any]
	created  bool
}

func NewCache(path Path, named ...string) Cache {
	var e std.Entity
	if len(named) == 0 {
		e = std.NewEntity[format.Default]()
	} else {
		e = std.NewEntityNamed(named[0])
	}
	i, _ := NewIdea(NewThought(make(map[string]any)), path)
	return Cache{
		Entity:   e,
		memories: i,
	}
}

func (cache Cache) sanityCheck() {
	if !cache.created {
		panic("a std.Cache must be created through std.NewCache")
	}
}

func (cache Cache) Reveal(path Path, code ...any) (any, error) {

}

func (cache Cache) Describe(revelation any, path Path, code ...any) error {

}

func (cache Cache) String() string {
	return cache.Entity.Named()
}
