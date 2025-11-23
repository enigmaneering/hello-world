package main

import (
	"image/color"
	"time"

	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution5/core/glitter"
	"git.ignitelabs.net/janos/core"
)

func main() {
	a := glitter.CreateWindowAt(640, 480, 55, 66, "Hello, glitter!", render)
	b := glitter.CreateWindowAt(640, 480, 77, 88, "Hello, glitter!", render)
	c := glitter.CreateWindowAt(640, 480, 99, 111, "Hello, glitter!", render)
	d := glitter.CreateWindowAt(640, 480, 122, 133, "Hello, glitter!", render)

	go a.Title("Window A")
	go b.Title("Window B")
	go c.Title("Window C")
	go d.Title("Window D")

	go func() {
		toggle := false
		for core.Alive() {
			time.Sleep(2 * time.Second)
			if toggle {
				glitter.FrameRate = 240
			} else {
				glitter.FrameRate = 60
			}
			toggle = !toggle
		}
	}()

	go func() {
		time.Sleep(time.Second)
		a.Maximize()
		time.Sleep(time.Second)
		a.Maximize()
	}()

	glitter.Orchestrate()
}

func render(frame glitter.Frame) {
	for y := uint(0); y < frame.Height; y++ {
		for x := uint(0); x < frame.Width; x++ {
			r := uint8((x + frame.Number) % 256)
			g := uint8((y + frame.Number) % 256)
			b := uint8(128)
			frame.Image.SetRGBA(int(x), int(y), color.RGBA{r, g, uint8(b), 255})
		}
	}
	frame.Present()
}
