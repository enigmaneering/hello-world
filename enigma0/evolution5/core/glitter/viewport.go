package glitter

import (
	"image"
	"sync"

	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/sys/rec"
	"github.com/veandco/go-sdl2/sdl"
)

type Viewport struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture
	width    int
	height   int
	id       uint32
	render   func(Frame)
	impulse  chan Frame

	mutex sync.Mutex
}

func NewViewport(width, height int, title string, render func(Frame)) (view *Viewport, err error) {
	rec.Verbosef(ModuleName, "creating viewport \"%s\"\n", title)

	// 0 - Create on the SDL thread

	Synchro.Send(func() {
		var window *sdl.Window
		window, err = sdl.CreateWindow(
			title,
			sdl.WINDOWPOS_CENTERED,
			sdl.WINDOWPOS_CENTERED,
			int32(width),
			int32(height),
			sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE,
		)
		if err != nil {
			return
		}

		actualW, actualH := window.GetSize()
		width = int(actualW)
		height = int(actualH)

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
			return
		}

		id, _ := window.GetID()
		view = &Viewport{
			window:   window,
			renderer: renderer,
			texture:  texture,
			width:    width,
			height:   height,
			id:       id,
			render:   render,
			impulse:  make(chan Frame),
		}

		mutex.Lock()
		viewports[id] = view
		framebuffers[id] = image.NewRGBA(image.Rect(0, 0, width, height))
		mutex.Unlock()
	})

	// 1 - Receive render impulses in a goroutine

	go func() {
		for core.Alive() {
			view.render(<-view.impulse)
		}
	}()

	return view, err
}

func (view *Viewport) Close() {
	view.mutex.Lock()
	defer view.mutex.Unlock()

	rec.Verbosef(ModuleName, "closing viewport \"%s\"\n", view.window.GetTitle())

	Synchro.Send(func() {
		if view.texture != nil {
			view.texture.Destroy()
		}
		if view.renderer != nil {
			view.renderer.Destroy()
		}
		if view.window != nil {
			view.window.Destroy()
		}
	})
}

func (view *Viewport) Resize(width, height int) (err error) {
	view.mutex.Lock()
	defer view.mutex.Unlock()

	Synchro.Send(func() {
		view.width = width
		view.height = height

		// Recreate texture with new dimensions
		if view.texture != nil {
			view.texture.Destroy()
		}

		var texture *sdl.Texture
		texture, err = view.renderer.CreateTexture(
			sdl.PIXELFORMAT_ABGR8888,
			sdl.TEXTUREACCESS_STREAMING,
			int32(width),
			int32(height),
		)
		if err != nil {
			return
		}
		view.texture = texture
	})
	return err
}

func (view *Viewport) Title(name ...string) string {
	var title string
	if len(name) == 0 {
		Synchro.Send(func() {
			title = view.window.GetTitle()
		})
		return title
	}
	Synchro.Send(func() {
		view.window.SetTitle(name[0])
		title = view.window.GetTitle()
	})
	return title
}

func (view *Viewport) Move(x, y uint32) {
	Synchro.Send(func() {
		view.window.SetPosition(int32(x), int32(y))
	})
}

func (view *Viewport) present(img *image.RGBA) {
	view.mutex.Lock()
	defer view.mutex.Unlock()

	Synchro.Send(func() {
		// Upload framebuffer to GPU texture
		view.texture.Update(nil, img.Pix, img.Stride)

		// Render texture to screen (no VSync - returns immediately)
		view.renderer.Clear()
		view.renderer.Copy(view.texture, nil, nil)
		view.renderer.Present() // Does NOT block
	})
}
