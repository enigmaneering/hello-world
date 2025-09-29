# `E0S3 - Lifecycles`
### `Alex Petz, Ignite Laboratories, September 2025`

---

### Neural Activation Lifecycles

So far, we've used actions and potentials to drive _looping_ execution.  However, as discussed at the beginning
of this enigma, there are _four_ unique ways to activate neural activity.  These four are broken into two categories -
long-running and one-shot.  Looping and stimulative activations are examples of long-running activations which self-recycle,
whereas triggered and impulse activations _don't._  Let's recap them for a moment

**0. Looping** - Reactivates when the last cycle has completed and the potential is high

**1. Stimulative** - Reactivates on every impulse the potential is high, regardless of prior execution

**2. Triggered** - Waits for the potential to go high, activates once, then decays

**3. Impulse** - Attempts to activate and then decays, regardless of activation

These four primitives are the foundation of all other neural activity!  For instance, let's say you want something
to _try_ and gracefully fail - an impulse makes more sense than a looping activation.  Whereas if you'd like to
serve web requests, an impulse would only serve _one_ before decaying.  More importantly, however, is the _mathematical_
comparisons between a _looping_ and _stimulative_ activation.  A loop provides a way to perform real time _integration_
of a value, while stimulation creates a means of _differential analysis._

An integral takes infinitesimally smaller approximations of a signal across time, while a differential acts directly on
each _point._  Because of this, point calculations of a signal can be made off of a high frequency stimulation while 
a looping integral calculates the "area under the curve" from each new batch of temporal data points

<picture>
<img alt="Temporal Fragment Shading" src="assets/E0S1D6 - Logical Activation.svg" width="500" style="display: block; margin-left: auto; margin-right: auto;">
</picture>

Here, we have _looping_ observation tracks which operate stimulatively through stepped activation, allowing each
_point_ to be integrated against from the coalesce loop concurrently.  This creates a neural _layering_ effect which I hope
to explore deeper as we progress.  This style of activation is _far_ more efficient than stimulatively activating
a single neural track, but that's where we'll start for now.

### Stimulative Activation

A stimulative activation visually executes like such

<picture>
<img alt="Stimulative Activation" src="assets/E0S1D4 - Stimulative Activation.svg" width="500" style="display: block; margin-left: auto; margin-right: auto;">
</picture>

The best way to demonstrate this is by stimulatively activating an endpoint which prints the moment it starts,
sleeps, and then prints the moment it ends.  This will interleave the print calls in an identifiable fashion

    func main() {
        c := std.NewCortex(std.RandomName())
        c.Frequency = 1 //hz
    
        c.Synapses() <- std.NewSynapse(lifecycle.Stimulative, "Print", Printer, nil)
    
        c.Spark()
        core.KeepAlive(time.Second * 5)
    }
    
    func Printer(imp *std.Impulse) {
        rec.Printf(imp.Bridge.String(), "START: %v\n", time.Now())
        time.Sleep(time.Second * 3)
        rec.Printf(imp.Bridge.String(), " STOP: %v\n", time.Now())
    }
---
    [Makya Beslier ⇝ Print] stimulating
    [Makya Beslier ⇝ Print] START: 2025-09-26 13:39:41.206825 -0700 PDT m=+1.026289751
    [Makya Beslier ⇝ Print] START: 2025-09-26 13:39:42.206683 -0700 PDT m=+2.026150376
    [Makya Beslier ⇝ Print] START: 2025-09-26 13:39:43.206496 -0700 PDT m=+3.025967001
    [Makya Beslier ⇝ Print] START: 2025-09-26 13:39:44.206697 -0700 PDT m=+4.026170585
    [Makya Beslier ⇝ Print] STOP: 2025-09-26 13:39:44.207635 -0700 PDT m=+4.027108293
    [Makya Beslier ⇝ Print] START: 2025-09-26 13:39:45.206659 -0700 PDT m=+5.026135543
    [Makya Beslier ⇝ Print] STOP: 2025-09-26 13:39:45.206765 -0700 PDT m=+5.026241543
    ...
    [core] instance shutting down
    ...
    [Makya Beslier ⇝ Print] decayed
    [Makya Beslier ⇝ Print] STOP: 2025-09-26 13:39:46.206617 -0700 PDT m=+6.026097293
    [Makya Beslier ⇝ Print] STOP: 2025-09-26 13:39:47.206946 -0700 PDT m=+7.026428335
    [Makya Beslier ⇝ Print] STOP: 2025-09-26 13:39:48.207084 -0700 PDT m=+8.026569793
    [core] instance shut down complete

Here you can see that the activation takes three seconds to complete, but with a stimulative lifecycle the
endpoint will be reactivated before the last completes.

### Triggered Activation

A triggered activation is a _one shot_ activation that waits for a condition to be met.  To demonstrate this,
let's set up an activation with a delayed potential

	now := time.Now().Add(time.Second * 5)
	c.Synapses() <- std.NewSynapse(lifecycle.Triggered, "Print", Printer, func(*std.Impulse) bool {
		if time.Now().After(now) {
			fmt.Printf("potential: high\n")
			return true
		}
		fmt.Printf("potential: low\n")
		return false
	})
---

    // 1hz    

    [Teresina Indge ⇝ Print] setting a trigger
    [Teresina Indge ⇝ Print] potential: low
    [Teresina Indge ⇝ Print] potential: low
    [Teresina Indge ⇝ Print] potential: low
    [Teresina Indge ⇝ Print] potential: low
    [Teresina Indge ⇝ Print] potential: high
    [Teresina Indge ⇝ Print] START: 2025-09-26 13:38:28.783709 -0700 PDT m=+5.026481042
    [Teresina Indge ⇝ Print] STOP: 2025-09-26 13:38:31.784501 -0700 PDT m=+8.027282459
    [Teresina Indge ⇝ Print] decayed

The activation waits for five seconds to pass before the potential goes high, allowing the action to fire.

### Impulsive Activation

For my networking friends, _impulsive_ activation is the UDP of neural activity.  Tactically, it's the _spray and 
pray_ of neural activations, but that's what can also make it effective in a pinch.  When bogged down behind unknowns
and obstacles, our impulses can often back up our decisions in unexpected ways.  Logically speaking you can think of
the higher level constructs as ways of intelligently driving the primal impulses.  For now, let's show how entropy
plays a part in impulsive activation

    // Generate a random boolean value on initialization
    var impulsePotential = rand.IntN(2) == 1

	c.Synapses() <- std.NewSynapse(lifecycle.Impulse, "Print", Printer, func(*std.Impulse) bool {
		return impulsePotential
	})
---
    // ~50% of the time
    [Vanda Nafziger ⇝ Print] impulsing
    [Vanda Nafziger ⇝ Print] potential: low
    [Vanda Nafziger ⇝ Print] decayed

    // ~50% of the time
    [Gelelemend Casella ⇝ Print] impulsing
    [Gelelemend Casella ⇝ Print] potential: high
    [Gelelemend Casella ⇝ Print] START: 2025-09-26 13:37:09.261751 -0700 PDT m=+1.019352418
    [Gelelemend Casella ⇝ Print] STOP: 2025-09-26 13:37:12.262417 -0700 PDT m=+4.020027835
    [Gelelemend Casella ⇝ Print] decayed

This is really useful when you need to fire something but have _absolutely zero control_ over its ability to execute.

In the next solution, we'll start putting together our first intelligent system using neural activations - a web server =) 

### Synaptic Decay

While it doesn't deserve an entire solution, it's important to note that the _impulse_ can be used to decay the synapse.
If an activation has deemed this beat to be the end of its neural lifecycle, you can set `impulse.Decay = true` to cause
the synapse to self-decay.