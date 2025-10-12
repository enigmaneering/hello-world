package std

import (
	"fmt"
	"time"

	"git.ignitelabs.net/janos/core/sys/given"
	"git.ignitelabs.net/janos/core/sys/given/format"
	"git.ignitelabs.net/janos/core/sys/id"
)

// Entity provides the basic foundation of all logical temporal components - it's simply a grouping of an ID, name, and genesis moment.
//
// NOTE: Entity's String function.
type Entity struct {
	id   uint64
	Name given.Name

	Genesis time.Time
}

// GetID returns the entity's fixed identifier value.
//
// NOTE: Identifiers are instance unique and cannot be assigned.
func (e Entity) GetID() uint64 {
	return e.id
}

// Named gets the entity's Name value as a string.
//
// NOTE: All given names come with implicit cultural descriptions and gender biases, which can be explored through Name directly.
// This specifically returns the given.Given type's Name value.
func (e Entity) Named() string {
	return e.Name.Name
}

// String returns the Entity's identifier and Name as "[ID](Name)"
func (e Entity) String() string {
	return fmt.Sprintf("[%d](%v)", e.id, e.Name.String())
}

// NewEntityNamed creates a new entity, assigns it a unique identifier, and gives it the provided name.
//
// See NewEntity and NewEntityNamed
func NewEntityNamed(name string) Entity {
	return NewEntity[format.Default](given.New(name))
}

// NewEntity creates a new entity, assigns it a unique identifier, gives it a random name, and sets its genesis moment.
//
// See NewEntity and NewEntityNamed
func NewEntity[T format.Format](name ...given.Name) Entity {
	i := id.Next()
	var g given.Name
	if len(name) > 0 {
		g = name[0]
	} else {
		g = given.Random[T]()
	}

	ne := Entity{
		id:      i,
		Name:    g,
		Genesis: time.Now(),
	}

	return ne
}
