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

// ModuleName provides the string identifier used by the `rec` package in logging.
var ModuleName = "glitter"

// FrameRate sets the global rate of presentation to the display.
var FrameRate = 60.0 // (in Hz)

// The Synchro is used to Send() code to execute on the SDL2 thread.
var Synchro = make(std.Synchro)

var windows = make(map[uint32]*Window)
var framebuffers = make(map[uint32]*image.RGBA)
var mutex = &sync.Mutex{}
var frameNumber = uint(0)

// Orchestrate begins the SDL2 system and facilitates the neural rendering of graphical contexts.
func Orchestrate() {
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
					for id, viewport := range windows {
						fb := framebuffers[id]
						presented := false
						presentLock := &sync.Mutex{}
						select {
						case viewport.impulse <- Frame{
							Image:  fb,
							Width:  uint(viewport.width),
							Height: uint(viewport.height),
							Delta:  delta,
							Number: frameNumber,
							Present: func() {
								presentLock.Lock()
								defer presentLock.Unlock()

								if !presented {
									presented = true
									viewport.present(fb)
								}
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
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					// Check which window was closed
					if viewport, ok := windows[e.WindowID]; ok {
						go viewport.Close()
						delete(windows, e.WindowID)
						delete(framebuffers, e.WindowID)

						if len(windows) == 0 {
							core.ShutdownNow()
						}
					}
				} else if e.Event == sdl.WINDOWEVENT_SIZE_CHANGED {
					if win, ok := windows[e.WindowID]; ok {
						w, h := win.window.GetSize()
						win.resize(uint(w), uint(h))
						framebuffers[e.WindowID] = image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
					}
				}
			default:
			}

			Synchro.Engage()
		}
	}
}
