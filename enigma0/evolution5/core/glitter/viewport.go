package glitter

import (
	"image"

	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution5/core/std"
)

// A Viewport is a std.Path to a glitter.Context.
type Viewport std.Path

type Context struct {
	*image.RGBA
	Width  uint
	Height uint
	Title  string

	// resizable, borderless, etc...
}

func (view *Viewport) Spark() *Viewport {
	// TODO: Create a window and connect it's output to the contextual image
	return view
}
