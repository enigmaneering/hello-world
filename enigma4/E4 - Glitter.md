# `E4 - Glitter`
### `Alex Petz, Ignite Laboratories, June 2025`

---

### Neural Rendering
I take a very firm stance on rendering: God doesn't use hardware acceleration!  Why?  Well, just look at our 
biology - it's a playground of _concurrent_ execution, not _parallel._  Things may _fire_ in a cluster, but
their operations do not _branch_ in lock-step. GPUs are wonderful if you'd like to perform the **exact same** 
operation for every _pixel_ in parallel, but why are we looking at the _pixel_ level?  Instead, I theorize 
that you can relegate the rendering of something to a _layered_ process which is distributed as a field of 
voxelized regions on the screen - entirely in software, using goroutines, and without any hardware acceleration =)

    "That's a bold strategy, Cotton! 
     ...Let's see if it pays off for him"

To begin, software rendering is actually quite straightforward using SDL2.  It's cross-platform, heavily
used in the industry, and only requires a _binding library_ to be called directly from Go.  Luckily, a
there's a wonderful group of people who love to do fun things called `veandco` who've lovingly created
the binding library we'll be using for this project - `github.com/veandco/go-sdl2`

Be warned, though - graphical programming is typically a _thread-strict_ environment!  Basic things like
window creation must be done from the _same_ thread that initializes SDL2, polls for events, and more.  Luckily,
most of the time _you_ don't care which thread your operations are executed on.  Because of that, I've written
what I call a `std.Synchro` which allows for non-blocking receiving of ad-hoc code while retaining blocking
execution from the calling thread.  This mitigates a majority of the threaded concerns that come with neural
execution through goroutines.

    package std
    
    import "sync"
    
    // Synchro represents a way to synchronize execution across threads.
    //
    // To send execution using a synchro, first create one using make.  Then, Engage the synchro
    // (non-blocking) from the thread you wish to execute on.  The calling thread can then Send
    // actions (blocking) to the other thread for intermittent execution.
    //
    //		 global -
    //	   var synchro = make(std.Synchro)
    //
    //		 main loop -
    //		  for ... {
    //	    ...
    //		   synchro.Engage()
    //		   ...
    //		  }
    //
    //		 sender -
    //	   synchro.Send(func() { ... })
    type Synchro chan *syncAction
    
    // syncAction represents a "waitable" action.
    type syncAction struct {
        sync.WaitGroup
        Action func()
    }
    
    // Send sends the provided action over the synchro channel and waits for it to be executed.
    func (s Synchro) Send(action func()) {
        syn := &syncAction{Action: action}
        syn.Add(1)
        s <- syn
        syn.Wait()
    }
    
    // Engage asynchronously handles -all- of the currently incoming actions on the Synchro channel before returning control.
    func (s Synchro) Engage() {
        for {
            select {
            case syn := <-s:
                syn.Action()
                syn.Done()
            default:
                return
            }
        }
    }
    
    // EngageOnce synchronously reads a single action on the Synchro channel before returning control.
    func (s Synchro) EngageOnce() {
        syn := <-s
        syn.Action()
        syn.Done()
    }

If you'd like to execute code directly on the SDL2 thread, you can do so from `glitter.Synchro`

    glitter.Synchro.Send(func() {
        fmt.Println(sdl.GetPlatform())
    })

### `glitter`
All rendering in JanOS is facilitated through the `glitter` package.  In glitter every window is considered to
be a _viewport_, rather than a _window._  This is because, in a neural architecture, these windows will often
represent viewports into the _same simulation._  Plus, we needed a unique type name, so the fact it denotes
such a difference is more of a happy little side effect =)

A viewport only needs _one_ function to operate - `render(glitter.Frame)` - and many viewports can share
the same render function.  Below is an example which creates four viewports that all render the same function,
which animates a rainbow tile pattern across the screen.

    func main() {
        go glitter.NewViewport(800, 600, "Window 0", render)
        go glitter.NewViewport(800, 600, "Window 1", render)
        go glitter.NewViewport(800, 600, "Window 2", render)
        go glitter.NewViewport(800, 600, "Window 3", render)
        glitter.Start()
    }
    
    func render(frame glitter.Frame) {
        for y := uint(0); y < frame.Height; y++ {
            for x := uint(0); x < frame.Width; x++ {
                r := uint8((x + frame.Number) % 256)
                g := uint8((y + frame.Number) % 256)
                b := uint8(128)
                frame.Image.SetRGBA(int(x), int(y), color.RGBA{r, g, b, 255})
            }
        }
    }

The first thing you'll notice is that `glitter` takes control of blocking your code and launching `core` - this
is a byproduct of how _macOS_ specifically handles rendering: SDL2 _**must**_ be launched from the main application
thread.  It's beyond my scope to try and mitigate that, so `glitter` acts as the application's central blocking 
point - but, thanks to the `std.Synchro` and the magic of channels, your rendering is happening in goroutines while
still being flipped to the screen on a regular interval from the SDL2 thread.

The _global_ framerate can be controlled through `glitter.FrameRate` - I emphasize _global_ because the concept of a
'global' framerate becomes a little absurd once you introduce neural timing on top of things =)