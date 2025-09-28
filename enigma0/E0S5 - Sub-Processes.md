# `E0S3 - Sub-Processes`
### `Alex Petz, Ignite Laboratories, September 2025`

---

### Neural Containers

The cortex provides a panic-safe way of launching your routines - but sometimes you'll
want to spark off _external_ processes!  Luckily, Go has a wonderful set of tools for doing
_exactly that_ - instead, let's try a challenge

    Create a program which can gracefully coordinate multiple sub-processes

For this example, I've pulled our neural _"Hello, World!"_ server example into the `neural.Net.HelloWorld`
method (to hopefully keep clutter out of our examples)

		cortex.Synapses() <- neural.Net.HelloWorld(lifecycle.Looping, std.RandomName(), os.Args[2])

Before we continue, though, it's time to discuss why there is so much 'naming' in this process!  _How else_ could you
track concurrent actions?  An integer identifier is so impersonal and hard to visually track, but 'Alice' is _immediately_
identifiable and traceable in the logs.  This will become _starkly_ apparent in a moment, but I promise there's rhyme
to my reason!  While you are _NOT_ required to use the human-legible name system, it's highly encouraged to create
_logically reasonable_ names.

Every core instance is given a name on creation, but you can define its _purpose_ through a description.  Any call to
`core.Describe` will set the global descriptor for the core instance, which immediately prints to the console (regardless
of verbosity)

    core.Describe("A Hello, World Server")
---
    [core] Averill Bohman is a "A 'Hello, World!' Server"

Any recordings by the `core` package will include the instance name next to `core` - allowing you to quickly distinguish
core activities in the logs.

### neural.Shell

Now, let's start exploring how to launch shell sub-processes!  It's quite simple, really

	cortex.Synapses() <- neural.Shell.SubProcess(lifecycle.Triggered, "Echo", []string{"echo", "Hello, World!"})
---
    [Anshu Witz ⇝ Echo] sparking sub-process 'echo'
    /bin/echo Hello, World!
    Hello, World!
    [Anshu Witz ⇝ Echo] sub process exited

Here we've fired off a single triggered activation that simply called the shell program `echo "Hello, World!"`, then
decayed.  The stdin/out/err of the sub-process is piped into the stdin/out/err of the host process, yielding the 
above output.  This allows you to _impulsively_ activate external processes using a channel - pretty neat!  Let's
explore it a little further and write a recursively activating program

    func main() {
        var cortex *std.Cortex
    
        if len(os.Args) > 1 && os.Args[1] == "server" {
            core.Describe("Sub-Process")
            cortex = std.NewCortex(std.RandomName())
    
            cortex.Synapses() <- neural.Net.HelloWorld(lifecycle.Looping, "Server", os.Args[2])
        } else {
            core.Describe("Multiplexer")
            cortex = std.NewCortex(std.RandomName())
    
            cortex.Synapses() <- neural.Shell.SubProcess(lifecycle.Looping, "sub process a", []string{"go", "run", "./main", "server", ":4242"}, func(imp *std.Impulse) {
                cortex.Impulse()
            })
            cortex.Synapses() <- neural.Shell.SubProcess(lifecycle.Looping, "sub process b", []string{"go", "run", "./main", "server", ":4243"}, func(imp *std.Impulse) {
                cortex.Impulse()
            })
            cortex.Synapses() <- neural.Shell.SubProcess(lifecycle.Looping, "sub process c", []string{"go", "run", "./main", "server", ":4244"}, func(imp *std.Impulse) {
                cortex.Impulse()
            })
        }
    
        cortex.Frequency = 1 //hz
        cortex.Mute()
        cortex.Spark()
        cortex.Impulse()
        core.KeepAlive()
    }

In this example, launching the program directly will spawn off three sub-processes - each launching a web server on its
provided port.  Depending on the branch taken, each will identify as a different kind of instance - but all will record
to the host processes stdout.  Let's take a look at the recorded output

    [core] Faustina Petruska is a "Multiplexer"
    [Hurst Littrick ⇝ sub process c] sparking sub-process 'go'
    [Hurst Littrick ⇝ sub process b] sparking sub-process 'go'
    /usr/local/go/bin/go run ./main server :4244
    /usr/local/go/bin/go run ./main server :4243
    [Hurst Littrick ⇝ sub process a] sparking sub-process 'go'
    /usr/local/go/bin/go run ./main server :4242
    [core] Araminta Maxwale is a "Sub-Process"
    [Januarius MacCallam ⇝ Server] neural server listening on :4242
    [core] Nemesio Amies is a "Sub-Process"
    [Suhad Adelman ⇝ Server] neural server listening on :4244
    [core] Digby Blasio is a "Sub-Process"
    [Helena Antonin ⇝ Server] neural server listening on :4243
    
    [core] Nemesio Amies instance shutting down
    [core] Nemesio Amies running 1 deferral
    [core] signing off — "Nemesio Amies, Sub-Process"
    
    [core] Digby Blasio instance shutting down
    [core] Digby Blasio running 1 deferral
    [core] signing off — "Digby Blasio, Sub-Process"
    
    [core] Faustina Petruska instance shutting down
    [core] Faustina Petruska running 1 deferral
    [core] signing off — "Faustina Petruska, Multiplexer"
    
    [core] Araminta Maxwale instance shutting down
    [core] Araminta Maxwale running 1 deferral
    [core] signing off — "Araminta Maxwale, Sub-Process"

Because of the neural naming system we can see that _Faustina_ created _Hurst,_ who in turn created _Araminta, Nemesio,_
and _Digby_ as servers.  In 45 lines of code (including import statements) we've launched off three web servers as a
standardized _process_ - allowing them to be coordinated as a cluster of activity.  In future solutions, we'll talk about
how to share a neural map between instances using a memory-mapped file =)
