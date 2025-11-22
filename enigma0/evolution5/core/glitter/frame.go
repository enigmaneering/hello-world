package glitter

import (
	"image"
	"time"
)

// A Frame is the canvas to draw a singular presentable frame upon.
type Frame struct {
	Image  *image.RGBA
	Width  uint
	Height uint

	Delta  time.Duration
	Number uint

	Present func()
}
