# `E0S2 - Is Everything Done?`
### `Alex Petz, Ignite Laboratories, September 2025`

---

### Neural Cleanup

The next important aspect of a neural architecture is the ability to _cleanup_ after oneself.  JanOS provides several
points you can introduce a cleanup operation at - from the instance level to the neural endpoint.  Instance level
deferral can be handled through the `core.Deferrals` channel

	core.Deferrals() <- func(wg *sync.WaitGroup) {
		time.Sleep(time.Second)
		wg.Done()
	}
	core.ShutdownNow()
---
    [core] instance shutting down
    [core] running 1 deferral
    [core] instance shut down complete

Calling `core.ShutdownNow()` will cause a cascading chain reaction of neural shutdowns - which can _also_ be deferred
against at multiple levels.  When neural activity ceases to cycle, its called _decaying._  To defer against the cortex 
decaying, you'd do the following

    func main() {
        c := std.NewCortex(std.RandomName())
        c.Frequency = 1 //hz

        c.Deferrals() <- func(wg *sync.WaitGroup) {
            fmt.Println("cortical deferral")
            time.Sleep(time.Second)
            wg.Done()
        }
        c.Spark()
    
        go core.Shutdown(time.Second * 5) // Delayed shutdown
        core.KeepAlive()
    }
---
    [core] instance shutting down in 5s
    
    [core] instance shutting down
    [core] running 1 deferral
    cortical deferral
    [core] instance shut down complete

In the above examples, the deferred functions are able to _hold open_ the shutdown process - however, for _synapses_
that doesn't apply.  A synapse simply can't _directly_ hold open the shutdown process, as it's transiently activated.
Because of that, you can either hold open the shutdown process by implementing your own mechanic using the provided 
deferral tools, or through passing a duration to `core.KeepAlive`.  The latter allows you to ensure all synapses can 
complete their cleanup functions within a _reasonable_ window of time.  A cleanup function can be given to a synapse
on creation through the optional final parameter

    func main() {
        c := std.NewCortex(std.RandomName())
        c.Frequency = 1 //hz
    
        c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Print", Printer, nil, Cleanup)

	    go core.Shutdown(time.Second * 3) // Delayed shutdown
    
        c.Spark()
        core.KeepAlive(time.Second * 5)
    }
    
    func Cleanup(imp *std.Impulse) {
        rec.Printf(imp.Bridge, "synapse cleaning up\n")
    }
    
    func Printer(imp *std.Impulse) {
        rec.Printf(imp.Bridge, "%v\n", imp.Timeline.CyclePeriod())
    }
---
    [Marelda Maytum ⇝ Print] 1.002303334s
    [Marelda Maytum ⇝ Print] 1.000355916s
    [Marelda Maytum ⇝ Print] 999.98525ms
    
    [core] instance shutting down
    [core] running 2 deferrals
    [core] holding open for 5s
    [Marelda Maytum ⇝ Print] cleaning up
    [core] instance shut down complete

### atlas

There's _so much more_ going on under the hood than meets the eye, which is where the _atlas_ configuration system
comes into play.  You may optionally provide an `atlas` file with your module, where you can set verbosity
and other things (like disabling the default preamble).  In the case of observing decaying, it's easier to understand
what's happening if we enable _verbose_ recording.

    NOTE: This may change from the JSON format in the future

    {
        "verbose": true,
        "printPreamble": false // Defaults to true if omitted
    }
---
    [core] created cortex 'Deorsa Readshaw'
    [core] creating synapse 'Print'
    [Deorsa Readshaw] sparking neural activity
    [core] instance shutting down in 3s
    [Deorsa Readshaw] wired axon to neural endpoint 'Print'
    [Deorsa Readshaw ⇝ Print] looping
    [Deorsa Readshaw ⇝ Print] 1.003269s
    [Deorsa Readshaw ⇝ Print] 999.879208ms
    
    [core] instance shutting down
    [core] running 2 deferrals
    [Deorsa Readshaw ⇝ Print] 999.923583ms
    [core] holding open for 5s
    [Deorsa Readshaw] cortex shutting down
    [Deorsa Readshaw] decayed
    [Deorsa Readshaw] cortex shut down complete
    [Deorsa Readshaw ⇝ Print] cleaning up
    [Deorsa Readshaw ⇝ Print] decayed
    [core] instance shut down complete

Every cortex, by default, assigns its own cleanup deferral to core on creation.  By viewing the verbose logs, you
immediately get visual confirmation that your entity has _finished_ decaying - an often useful sanity check!  In the
land of debugging concurrent execution, knowing that an entity _hasn't_ decayed often reveals deadlocks and stalled conditions.

The atlas is unique to your _module_ - not your main file - and must live next to your `go.mod` file.  