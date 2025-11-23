package glitter

import (
	"image"
	"time"
)

// A Frame provides the canvas to draw a single atomic frame upon.
//
// When your drawing is complete, please call Present.
type Frame struct {
	// Image can be freely drawn upon to your wish.
	Image *image.RGBA

	// Width represents the width of this particular frame.
	Width uint

	// Height represents the height of this particular frame.
	Height uint

	// Delta represents the amount of time that has passed since the last tick of the glitter orchestration clock.
	//
	// NOTE: A -single- frame will be requested for each window per every cycle of the "clock."  This means you (the
	// window) can consider this value to be a reliable 'tick' interval to perform rendering calculations against.
	Delta time.Duration

	// Number is an ever incrementing frame number.
	Number uint

	// Present is what ultimately displays your image to the screen. It, by design, can
	// only be called once - as a frame represents a -single- atomic rendering operation.
	Present func()
}
