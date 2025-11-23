package glitter

import "C"
import (
	"image"
	"strconv"
	"sync"
	"time"

	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/sys/rec"
	"github.com/veandco/go-sdl2/sdl"
)

// A Window is the actual structure that manages an SDL2 window.
type Window struct {
	Error       error
	window      *sdl.Window
	renderer    *sdl.Renderer
	texture     *sdl.Texture
	width       uint
	height      uint
	id          uint32
	render      func(Frame)
	impulse     chan Frame
	destroyed   bool
	initialized bool

	mutex sync.Mutex
}

// String returns the window's ID as a string
func (win Window) String() string {
	return strconv.Itoa(int(win.id))
}

// CreateWindow creates and returns a handle to a hosting operating system's window, which renders the provided function at the global FrameRate.
func CreateWindow(width, height uint, title string, render func(Frame)) (win *Window) {
	return CreateWindowAt(width, height, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, title, render)
}

// CreateWindowAt creates and returns a handle to a hosting operating system's window explicitly positioned at creation, which renders the provided function at the global FrameRate.
func CreateWindowAt(width, height uint, x, y uint, title string, render func(Frame)) (win *Window) {
	win = &Window{}
	go func() {
		var err error

		// 0 - Create on the SDL thread

		Synchro.Send(func() {
			var window *sdl.Window
			window, err = sdl.CreateWindow(
				title,
				int32(x),
				int32(y),
				int32(width),
				int32(height),
				sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE,
			)
			if err != nil {
				return
			}

			actualW, actualH := window.GetSize()
			width = uint(actualW)
			height = uint(actualH)

			var renderer *sdl.Renderer
			renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
			if err != nil {
				rec.Verbosef(ModuleName, "hardware acceleration with VSync failed, trying without VSync\n")

				// Try hardware without VSync first
				renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
				if err != nil {
					rec.Verbosef(ModuleName, "hardware acceleration failed completely, using software renderer\n")

					// Last resort: software renderer (VSync likely won't work anyway)
					renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
					if err != nil {
						window.Destroy()
						win.destroyed = true
						panic("glitter cannot run on this hardware")
						return
					}
				}
			}

			var texture *sdl.Texture
			texture, err = renderer.CreateTexture(
				sdl.PIXELFORMAT_ABGR8888,
				sdl.TEXTUREACCESS_STREAMING,
				int32(width),
				int32(height),
			)
			if err != nil {
				renderer.Destroy()
				window.Destroy()
				win.destroyed = true
				return
			}

			id, _ := window.GetID()
			win.window = window
			win.renderer = renderer
			win.texture = texture
			win.width = width
			win.height = height
			win.id = id
			win.render = render
			win.impulse = make(chan Frame)
			win.initialized = true
			rec.Verbosef(ModuleName, "created window [%s]\n", win)

			mutex.Lock()
			windows[id] = win
			framebuffers[id] = image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
			mutex.Unlock()
		})
		win.Error = err

		// 1 - Receive render impulses in a goroutine
		if win.Error == nil {
			go func() {
				for core.Alive() {
					win.render(<-win.impulse)
				}
			}()
		} else {
			win.Close()
		}
	}()
	return win
}

// Close destroys the window (in SDL2 terminology) and prevents further operations from completing against the window.
func (win *Window) Close() {
	win.mutex.Lock()
	defer win.mutex.Unlock()

	rec.Verbosef(ModuleName, "closing window [%s]\n", win)

	if win.texture != nil {
		win.texture.Destroy()
	}
	if win.renderer != nil {
		win.renderer.Destroy()
	}
	if win.window != nil {
		win.window.Destroy()
	}
	win.destroyed = true
}

// Resize attempts to resize the window.
//
// NOTE: This may cause a texture read/write assertion error in your console - that's because a frame might be
// mid-drawing while the underlying dimensions change.  It's not a problem, please ignore the message.
func (win *Window) Resize(width, height uint) (err error) {
	Synchro.Send(func() {
		win.resize(width, height)
	})
	return err
}

func (win *Window) resize(width, height uint) (err error) {
	if !win.sanityCheck() {
		return
	}
	rec.Verbosef(ModuleName, "resizing window [%s] to %dx%d\n", win, width, height)

	win.mutex.Lock()
	defer win.mutex.Unlock()

	win.width = width
	win.height = height

	// Recreate texture with new dimensions
	if win.texture != nil {
		win.texture.Destroy()
	}

	var texture *sdl.Texture
	texture, err = win.renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		int32(width),
		int32(height),
	)
	if err != nil {
		return
	}
	win.texture = texture
	return err
}

// Focus attempts to bring the window into focus, raising it above other windows, and returns whether it was able to do so.
//
// NOTE: The success of this is subject to many factors - see https://wiki.libsdl.org/SDL3/SDL_RaiseWindow
func (win *Window) Focus() bool {
	if !win.sanityCheck() {
		return false
	}

	var focused bool
	Synchro.Send(func() {
		win.window.Raise()
		flags := win.window.GetFlags()
		focused = (flags & sdl.WINDOW_INPUT_FOCUS) != 0
	})
	return focused
}

func (win *Window) Maximize() bool {
	if !win.sanityCheck() {
		return false
	}

	var maximized bool
	Synchro.Send(func() {
		win.window.Maximize()
		flags := win.window.GetFlags()
		maximized = (flags & sdl.WINDOW_MAXIMIZED) != 0
	})
	return maximized
}

func (win *Window) Minimize() bool {
	if !win.sanityCheck() {
		return false
	}

	var minimized bool
	Synchro.Send(func() {
		win.window.Minimize()
		flags := win.window.GetFlags()
		minimized = (flags & sdl.WINDOW_MINIMIZED) != 0
	})
	return minimized
}

// Title sets and/or gets the window's title.
//
// NOTE: To only get the title, please provide no parameters.
func (win *Window) Title(name ...string) string {
	if !win.sanityCheck() {
		return ""
	}

	var title string
	if len(name) == 0 {
		Synchro.Send(func() {
			title = win.window.GetTitle()
		})
		return title
	}
	rec.Verbosef(ModuleName, "setting window [%s] title to \"%s\"\n", win, name[0])
	Synchro.Send(func() {
		win.window.SetTitle(name[0])
		title = win.window.GetTitle()
	})
	return title
}

// SetPosition attempts to move the window to the desired coordinates.
func (win *Window) SetPosition(x, y uint32) {
	if !win.sanityCheck() {
		return
	}
	rec.Verbosef(ModuleName, "moving window [%s] to (%d, %d)\n", win, x, y)

	Synchro.Send(func() {
		win.window.SetPosition(int32(x), int32(y))
	})
}

// GetPosition gets the window's current coordinates.
func (win *Window) GetPosition() (x, y uint32) {
	if !win.sanityCheck() {
		return 0, 0
	}

	Synchro.Send(func() {
		xI, yI := win.window.GetPosition()
		x = uint32(xI)
		y = uint32(yI)
	})
	return x, y
}

func (win *Window) present(img *image.RGBA) {
	if !win.sanityCheck() {
		return
	}

	win.mutex.Lock()
	defer win.mutex.Unlock()

	Synchro.Send(func() {
		win.texture.Update(nil, img.Pix, img.Stride)

		win.renderer.Clear()
		win.renderer.Copy(win.texture, nil, nil)
		win.renderer.Present()
	})
}

func (win *Window) sanityCheck() bool {
	if win.destroyed {
		return false
	}
	for !win.initialized {
		time.Sleep(time.Millisecond)
	}
	return true
}
