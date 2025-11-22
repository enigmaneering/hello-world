package main

import (
	"image/color"
	"time"

	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution5/core/glitter"
)

func main() {
	go func() {
		a, _ := glitter.NewViewport(800, 600, "Window 0", render)
		glitter.NewViewport(800, 600, "Window 1", render)
		glitter.NewViewport(800, 600, "Window 2", render)
		glitter.NewViewport(800, 600, "Window 3", render)

		time.Sleep(time.Second * 5)
		for i := uint32(0); i < 5; i++ {
			for ii := uint32(0); ii < 5; ii++ {
				a.Move(i, ii)
			}
		}
	}()

	glitter.Start()
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
