# `E4 - Glitter`
### `Alex Petz, Ignite Laboratories, June 2025`

---

### Neural Rendering
I take a very firm stance on rendering: God doesn't use hardware acceleration!  Why?  Well, just look at our 
biology - it's a playground of _concurrent_ execution, not _parallel._  Things may _fire_ in a cluster, but
their operations do not _branch_ in lock-step as is expected with _parallel_ execution. GPUs are absolutely
wonderful and amazing pieces of machinery, but their power is in empowering _math_ - rendering has just been
a means to that end.  Instead, I theorize that you can relegate the rendering of something to a _layered_ and
_temporal_ process which is distributed as a field of voxelized regions on the screen - entirely in software, 
using goroutines, and without any hardware acceleration =)

    "That's a bold strategy, Cotton! 
     ...Let's see if it pays off for him"

To begin, software rendering is actually quite straightforward using SDL2.  It's cross-platform, heavily
used in the industry, and requires only a _binding library_ to be called directly from Go.  Luckily, a
there's a wonderful group of people who love to do fun things named `veandco` who've lovingly created
the binding library we'll be using for this project - `github.com/veandco/go-sdl2`

Be warned, though - graphical programming is typically a _thread-strict_ environment!  Even basic things like
window creation must be done from the _same thread_ that initializes SDL2, polls for events, and more.  Luckily,
most of the time _you_ don't care which thread your operations are executed on.  Because of that, I've written
what I call a `std.Synchro` which allows for "intermittent" remote execution of ad-hoc code while blocking the calling 
thread.  This mitigates a majority of the threaded concerns that come with neural execution through goroutines.

Think of it like a synchro in your vehicle, reactively spinning up to speed to make sure that both ends meet at the appropriate
~~speed~~ line number across time.

    package std
    
    import "sync"
    
    // Synchro represents a way to synchronize execution across threads.
    //
    // To send execution using a synchro, first create one using make - then, Engage the synchro
    // from the thread you wish to execute on.  The calling thread can then Send blocking actions
    // which the synchronizable thread should "intermittently" execute.
    //
    //    global -
    //  var synchro = make(std.Synchro)
    //
    //    main loop -
    //  for ... {
    //    ...
    //    synchro.Engage()
    //    ...
    //  }
    //
    //    sender -
    //  synchro.Send(func() { ... })
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
    //
    // NOTE: If you'd like to process ALL available messages in a single engagement, rather than one, please provide 'true' to 'processAll'
    func (s Synchro) Engage(processAll ...bool) {
        all := len(processAll) > 0 && processAll[0]
        for {
            select {
                case syn := <-s:
                    syn.Action()
                    syn.Done()
                    if !all {
                        return
                    }
                default:
                    return
            }
        }
    }

    // EngageBlocking synchronously handles the currently incoming actions on the Synchro channel before returning control.
    //
    // NOTE: If you'd like to process ALL available messages in a single engagement, rather than one, please provide 'true' to 'processAll'
    func (s Synchro) EngageBlocking(processAll ...bool) {
        all := len(processAll) > 0 && processAll[0]
        for {
            select {
                case syn := <-s:
                    syn.Action()
                    syn.Done()
                    if !all {
                        return
                    }
            }
        }
    }

If you'd like to execute code directly on the main SDL2 thread, you can send it over via the `glitter.Synchro`

    glitter.Synchro.Send(func() {
        fmt.Println(sdl.GetPlatform())
    })

### `glitter`
All rendering in JanOS is facilitated through the `glitter` package.  We begin with a `glitter.Window` - which,
as you would expect, is how you create an abstract _window_ in any operating system.  What's different from an
`sdl.Window` is that a `glitter.Window` performs all of the necessary wiring to present you with a simple
`image.RGBA` you can _immediately_ begin to draw upon.  Creating a `glitter.Window` is a breeze

    func main() {
        go glitter.CreateWindow(640, 480, "Hello, glitter!", render)

        glitter.Orchestrate() // This facilitates neural orchestration of graphical contexts
    }

    // This is a generic tiling rainbow pattern which animates with the frame.Number
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

The call to `glitter.Orchestrate()` is a blocking operation which integrates `core`, SDL2, and the mind
map endpoints needed to facilitate neural rendering.

Here, we've provided a `render` function that operates against a `glitter.Frame`, which represents
the canvas for each rendering operation to occur against.  A window has the standard movement and naming
operations that come with it, but as these operations are _blocking_ you'll need to ensure that 
`glitter.Orchestrate()` can still be reached.  For example, let's create four duplicate windows and 
immediately rename them all 

    func main() {
        a := glitter.CreateWindow(640, 480, "Hello, glitter!", render)
        b := glitter.CreateWindow(640, 480, "Hello, glitter!", render)
        c := glitter.CreateWindow(640, 480, "Hello, glitter!", render)
        d := glitter.CreateWindow(640, 480, "Hello, glitter!", render)
    
        // These operations must be performed in a goroutine
        go a.Title("Window A") // NOTE: Title is a setter and/or getter =)
        go b.Title("Window B")
        go c.Title("Window C")
        go d.Title("Window D")
    
        glitter.Orchestrate()
    }
    
    // This is a generic tiling rainbow pattern which animates with the frame.Number
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

A _window_ is only a means to an end - you genuinely shouldn't have to interact at the _window_ level in
the first place!  I wrote glitter because in 2025 our children _**still**_ have to learn about window compositing,
v-sync, hardware acceleration, and think in _single threaded_ mindsets if they'd like to just _**draw something**_
in code!

What gives!?  Shouldn't that be one of the most primitive tools they have access to by now???

Our kids should be able to think in terms of _vector graphics_ without having to understand the tesselation needed
to render it out to the screen!

_That's how we entice them to **eventually** learn about shader programs, hardware acceleration, and advanced rendering techniques_ =)

A window gives the absolute _bare minimum_ necessary to start drawing in as few of lines of code as possible - but there's another way!

### Viewports

A `glitter.Viewport` is what bridges a window to the mind map.  A viewport exists as an outlet to 'peek' into whatever
is found at the end of a lick.  At runtime, the viewport will evaluate the lick and determine how to create a window 
displaying whatever it is that it found

    glitter.Viewport(std.Path{"a", "Path", 2, "something"}).Spark()

While that's cool, what's even _cooler_ are the standardized viewports you can spark into creation =)

    glitter.Spark.Trace(*value, dimensions)

Here, a window will be created which plots a value across time using the provided dimensional constraints - essentially
creating a semantic _oscilloscope_ =)