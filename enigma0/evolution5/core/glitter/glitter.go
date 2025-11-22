package glitter

import (
	"image"
	"sync"
	"time"

	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution5/core/std"
	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/sys/rec"
	"git.ignitelabs.net/janos/core/sys/when"
	"github.com/veandco/go-sdl2/sdl"
)

var ModuleName = "glitter"

var FrameRate = 60.0 // hz

// Synchro can be used to execute code on the SDL2 rendering thread.
var Synchro = make(std.Synchro)

var viewports = make(map[uint32]*Viewport)
var framebuffers = make(map[uint32]*image.RGBA)
var mutex = &sync.Mutex{}
var frameNumber = uint(0)

func Start() {
	rec.Verbosef(ModuleName, "initializing SDL2\n")

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		rec.Fatalf(ModuleName, err.Error())
	}
	defer sdl.Quit()

	last := time.Now()
	next := last.Add(when.HertzToDuration(FrameRate))

	go func() {
		for core.Alive() {
			if time.Now().After(next) {
				n := time.Now()
				delta := n.Sub(last)
				last = n
				next = last.Add(when.HertzToDuration(FrameRate))

				Synchro.Send(func() {
					mutex.Lock()
					for id, viewport := range viewports {
						fb := framebuffers[id]
						select {
						case viewport.impulse <- Frame{
							Image:  fb,
							Width:  uint(viewport.width),
							Height: uint(viewport.height),
							Delta:  delta,
							Number: frameNumber,
							Present: func() {
								viewport.present(fb)
							},
						}:
						default:
						}
					}
					mutex.Unlock()

					frameNumber++
				})
			}
		}
	}()

	for core.Alive() {
		Synchro.Engage()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				core.ShutdownNow()
			case *sdl.RenderEvent:

			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					// Check which window was closed
					if viewport, ok := viewports[e.WindowID]; ok {
						viewport.Close()
						delete(viewports, e.WindowID)
						delete(framebuffers, e.WindowID)

						if len(viewports) == 0 {
							core.ShutdownNow()
						}
					}
				} else if e.Event == sdl.WINDOWEVENT_SIZE_CHANGED {
					if viewport, ok := viewports[e.WindowID]; ok {
						w, h := viewport.window.GetSize()
						rec.Verbosef(ModuleName, "resizing viewport \"%s\" to %dx%d\n", viewport.window.GetTitle(), w, h)
						viewport.Resize(int(w), int(h))
						framebuffers[e.WindowID] = image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
					}
				}
			}
		}
	}
}
