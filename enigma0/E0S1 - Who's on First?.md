# `E0S1 - Who's on First?`
### `Alex Petz, Ignite Laboratories, September 2025`

---

### Mixed Timing

In the last example we _named_ our neural endpoints, but how do we know who is calling who?  In a neural
architecture, many different cortices can work in tandem to drive impulses.  Because of this, cortices
are given a name at creation - often through `std.RandomName()` (though there's nothing stopping you from
naming your cortex 'motor control' or 'proprioception')

A synapse, however, is a _bridge_ between a cortex and neuron that exists only in the scope of a neural activation.
After all, a synapse is a _function_ - not a type. Thus, it makes more sense to describe a synapse by the bridge it 
forms.  Let's take a look at what a bridge looks like

    func main() {
        c := std.NewCortex(std.RandomName())
        c.Frequency = 1 //hz
    
        rec.Printf(c.Named(), "Hello, World!\n")
    
        c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Print", Printer, nil)
    
        c.Spark()
        core.KeepAlive()
    }
    
    func Printer(imp *std.Impulse) {
        rec.Printf(imp.Bridge.String(), "%v\n", imp.Timeline.CyclePeriod())
    }
---
    [Shafiq Dedman] Hello, World!
    [Shafiq Dedman ⇝ Print] 1.003896542s
    [Shafiq Dedman ⇝ Print] 1.000341042s
    [Shafiq Dedman ⇝ Print] 1.016742625s

The _bridge_ - analogous to an axon - is a _named_ connection between a cortex and synapse.  It's only visible in
the recorded output and only _indirectly_ controllable (which we'll talk about later).  When creating a synapse 
in this way, you technically create an implicit neuron.  If you'd like to create many synaptic connections to a single 
neuron, you can!  The reason to do so is when you have multiple cortices driving the _same_ neural endpoint

    func main() {
        c1 := std.NewCortex(std.RandomName())
        c1.Frequency = 1 //hz
    
        rec.Printf(c1.Named(), "Hello, World!\n")
    
        c2 := std.NewCortex(std.RandomName())
        c2.Frequency = 1 //hz
    
        rec.Printf(c2.Named(), "Hello, World!\n")
    
        n := std.NewNeuron("Print", Printer, nil)
        c1.Synapses() <- std.NewSynapseFromNeural(lifecycle.Looping, n)
        c2.Synapses() <- std.NewSynapseFromNeural(lifecycle.Looping, n)
    
        c1.Spark()
        c2.Spark()
        core.KeepAlive()
    }
    
    func Printer(imp *std.Impulse) {
        rec.Printf(imp.Bridge.String(), "%v\n", imp.Timeline.CyclePeriod())
    }
---
    [Leon Adlem] Hello, World!
    [Alva Riccio] Hello, World!
    [Alva Riccio ⇝ Print] 1.003297375s
    [Leon Adlem ⇝ Print] 1.003302292s
    [Alva Riccio ⇝ Print] 999.873042ms
    [Leon Adlem ⇝ Print] 999.881375ms
    [Alva Riccio ⇝ Print] 1.000007375s
    [Leon Adlem ⇝ Print] 1.000005916s

Let's not get too far ahead, though, and stick to a single cortex for now.  Instead, let's create multiple _neural
endpoints_ that activate the same function at _different_ rates.  To do so, we'll create implicit neurons through
simple synapses again
    
    func main() {
        c := std.NewCortex(std.RandomName())
        c.Frequency = 60 //hz
    
        rec.Printf(c.Named(), "Hello, World!\n")
    
        c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Print A", Printer, when.Frequency[float64](1))
        c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Print B", Printer, when.Frequency[float64](0.7))
        c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Print C", Printer, when.Frequency[float64](0.4))
    
        c.Spark()
        core.KeepAlive()
    }
    
    func Printer(imp *std.Impulse) {
        rec.Printf(imp.Bridge.String(), "%v\n", imp.Timeline.CyclePeriod())
    }
---
    [Kiara Bridat] Hello, World!
    [Kiara Bridat ⇝ Print A] 1.002363333s
    [Kiara Bridat ⇝ Print B] 1.4353665s
    [Kiara Bridat ⇝ Print A] 1.000167375s
    [Kiara Bridat ⇝ Print C] 2.502636375s
    [Kiara Bridat ⇝ Print B] 1.433120833s
    [Kiara Bridat ⇝ Print A] 1.000188917s
    [Kiara Bridat ⇝ Print A] 1.016833917s
    [Kiara Bridat ⇝ Print B] 1.434447125s
    [Kiara Bridat ⇝ Print C] 2.516854958s

Here we've begun to use the `when` package - which provides a bunch of ways to create type-agnostic temporal potential 
functions.  More importantly, you can create potentials off of _references!_

    func main() {
        c := std.NewCortex(std.RandomName())
        c.Frequency = 60 //hz
    
        rec.Printf(c.Named(), "Hello, World!\n")
    
        c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Print", Printer, when.FrequencyRef(&frequency))
    
        c.Spark()
        core.KeepAlive()
    }
    
    var frequency = 1.0
    var toggle = true
    
    func Printer(imp *std.Impulse) {
	    // NOTE: This scheme doesn't guarantee an even timing interval =)

        if toggle {
            frequency = 0.5
        } else {
            frequency = 1.0
        }
        toggle = !toggle

        rec.Printf(imp.Bridge.String(), "%v\n", imp.Timeline.CyclePeriod())
    }
---
    [Jeannine Brunelleschi] Hello, World!
    [Jeannine Brunelleschi ⇝ Print] 1.001791916s
    [Jeannine Brunelleschi ⇝ Print] 2.000328s
    [Jeannine Brunelleschi ⇝ Print] 1.000002209s
    [Jeannine Brunelleschi ⇝ Print] 2.000018041s
    [Jeannine Brunelleschi ⇝ Print] 1.000030584s
    [Jeannine Brunelleschi ⇝ Print] 2.000053166s

_This_ is what makes neural activation so powerful!  The neurons can control _their own_ activation intelligently.
In the next solution, we'll start talking about how neurons can _decay_ and perform self-cleanup =)