# `E0S0 - How Long has Passed?`
### `Alex Petz, Ignite Laboratories, September 2025`

---

### Neural Timekeeping

Let's pose ourselves a challenge:

    Fire a function at a controlable interval of time which prints out its moment of activation

Now, let's take the naïve approach:

    for true {
        fmt.Println(time.Now())
        time.Sleep(time.Second)
    }

Here's the output:

    2025-09-24 17:33:10.09763 -0700 PDT m=+0.000047043
    2025-09-24 17:33:11.098818 -0700 PDT m=+1.001235918
    2025-09-24 17:33:12.099735 -0700 PDT m=+2.002154168
    2025-09-24 17:33:13.100852 -0700 PDT m=+3.003272001
    2025-09-24 17:33:14.101975 -0700 PDT m=+4.004395501
    2025-09-24 17:33:15.103081 -0700 PDT m=+5.005501709
    2025-09-24 17:33:16.103213 -0700 PDT m=+6.005635334
    2025-09-24 17:33:17.103742 -0700 PDT m=+7.006164709
    2025-09-24 17:33:18.10403 -0700 PDT m=+8.006453334
    2025-09-24 17:33:19.105107 -0700 PDT m=+9.007530876
    2025-09-24 17:33:20.105987 -0700 PDT m=+10.008412293
    2025-09-24 17:33:21.106994 -0700 PDT m=+11.009419209
    2025-09-24 17:33:22.107249 -0700 PDT m=+12.009674959
    2025-09-24 17:33:23.108355 -0700 PDT m=+13.010781959
    2025-09-24 17:33:24.10947 -0700 PDT m=+14.011897668
    2025-09-24 17:33:25.110578 -0700 PDT m=+15.013007126
    2025-09-24 17:33:26.111716 -0700 PDT m=+16.014145793
    2025-09-24 17:33:27.112624 -0700 PDT m=+17.015054043
    2025-09-24 17:33:28.113715 -0700 PDT m=+18.016145376
    2025-09-24 17:33:29.114882 -0700 PDT m=+19.017313918

So far so good, right?  Let's "scale" in by speeding up to 16ms (~60hz) and _calculate_ the refractory period

	last := time.Now()
	for true {
		now := time.Now()
		fmt.Println(now.Sub(last))
		last = now
		time.Sleep(time.Millisecond*16) // ~60hz
	}

Here's the output:

    17.038667ms
    17.028541ms
    17.040625ms
    17.013667ms
    17.00875ms
    17.019875ms
    17.017667ms
    17.014083ms

_That's not 16ms!_  In fact, the more steps you add to the loop, the further you deviate from your target frequency.
To solve this, you can implement a simple feedback loop

    interval := time.Millisecond * 16

	last := time.Now()
	var adjustment time.Duration
	for true {
        expected := last.Add(interval).Sub(time.Now().Add(adjustment))
        time.Sleep(expected)
        observed := time.Since(last)
        adjustment = observed - expected

        // Calculate from the same moment
        now := time.Now() 
        fmt.Println(now().Sub(last))
        last = now()
	}

Now, the output fluctuates _around_ our desired interval:

    17.012ms
    15.690042ms
    15.682625ms
    16.62ms
    15.997541ms
    15.995541ms
    16.000416ms
    16.007375ms

Wonderful!  We've solved our issue, right?  Well - yes, but that's a _mouthful_ of code to drop into every
loop.  It would be a _lot_ simpler if we could simply _ask_ a function to be called with such timing


    func main() {
        c := std.NewCortex(std.RandomName())
        c.Frequency = 60 //hz (16.‾6ms)
        
        c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Printer", Printer, nil)
    
        c.Spark()
        core.KeepAlive()
    }
        
    func Printer(imp *std.Impulse) {
        fmt.Println(imp.Timeline.CyclePeriod())
    }

The output of this, as you would expect, is just as accurate as our handwritten feedback loop above. While this 
doesn't appear any smaller, the only line that actually creates a neural loop is the `c.Synapses()` line.  Everything
else is just some basic configuration.  More importantly, the `Printer` function has _zero_ timing clutter in its 
footprint, aside from the function signature!  

So, let's multiplex a few loops:

    func main() {
        c := std.NewCortex(std.RandomName())
        c.Frequency = 60 //hz (16.‾6ms)
    
        c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Cycle Printer", PrintCycle, nil)
        c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Refractory Printer", PrintRefractory, nil)
        c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Response Time Printer", PrintResponse, nil)
    
        c.Spark()
        core.KeepAlive()
    }
    
    func PrintCycle(imp *std.Impulse) {
        fmt.Println("[cycle] " + imp.Timeline.CyclePeriod().String())
    }
    
    func PrintRefractory(imp *std.Impulse) {
        fmt.Println("[refraction] " + imp.Timeline.RefractoryPeriod().String())
    }
    
    func PrintResponse(imp *std.Impulse) {
        fmt.Println("[response] " + imp.Timeline.ResponseTime().String())
    }

Our output now contains the printed results of all three functions _independently_ being called in unison across
goroutines

    [response] 41ns
    [refraction] 19.649958ms
    [cycle] 19.646625ms
    [cycle] 16.620833ms
    [refraction] 16.610875ms
    [response] 125ns
    [cycle] 16.686208ms
    [refraction] 16.687083ms
    [response] 83ns
    [response] 208ns
    [cycle] 16.672667ms
    [refraction] 16.658416ms
    [refraction] 16.646042ms
    [response] 42ns
    [cycle] 16.678ms
    [cycle] 16.669ms

The first thing you'll notice is there's _no_ guarantee in the order of execution!  That can be solved by adding a
_potential function_ to our synapses.  Currently, we've passed `nil` as our potential (which equates to `when.Always()`)
but we can use any potential we'd like.  The `when` package provides a handful of common neural potentials, but in
this case we need a `when.StepMaker` - which creates a pair of functions we can use to control the neural activations
    
    var Step func()
    
    func main() {
        c := std.NewCortex(std.RandomName())
        c.Frequency = 60 //hz
    
        makePotential, step := when.StepMaker(3) // There are 3 neural endpoints to step between
        Step = step
    
        c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Cycle Printer", PrintCycle, makePotential(0))
        c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Refractory Printer", PrintRefractory, makePotential(1))
        c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Response Time Printer", PrintResponse, makePotential(2))
    
        c.Spark()
        core.KeepAlive()
    }
    
    func PrintCycle(imp *std.Impulse) {
        fmt.Printf("%v [cycle] %v\n", imp.Timeline.CyclePeriod().String(), imp.Timeline.CyclePeriod().String())
        Step()
    }
    
    func PrintRefractory(imp *std.Impulse) {
        fmt.Printf("%v [refraction] %v\n", imp.Timeline.CyclePeriod().String(), imp.Timeline.RefractoryPeriod().String())
        Step()
    }
    
    func PrintResponse(imp *std.Impulse) {
        fmt.Printf("%v [response] %v\n", imp.Timeline.CyclePeriod().String(), imp.Timeline.ResponseTime().String())
        Step()
    }

I've added the cycle period to the left so you can see the frequency of each function's activation: 20hz.  That's because
we're stepping between three endpoints in an infinite loop

    49.555916ms [cycle] 49.555916ms
    49.810417ms [refraction] 49.784167ms
    49.178834ms [response] 334ns
    50.846709ms [cycle] 50.846709ms
    50.039458ms [refraction] 50.021542ms
    50.874583ms [response] 167ns
    49.984125ms [cycle] 49.984125ms
    49.215375ms [refraction] 49.185208ms
    49.793417ms [response] 16.709µs

The potential function, however, can be anything you'd desire - let's put a little faux delay in the last activation

	c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Cycle Printer", PrintCycle, makePotential(0))
	c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Refractory Printer", PrintRefractory, makePotential(1))

	i := 0
	final := makePotential(2)
	c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Response Time Printer", PrintResponse, func(imp *std.Impulse) bool {
		if final(imp) {
			i++
			if i >= 128 {
				i = 0
				return true
			}
		}
		return false
	})

Here the final activation requires 128 impulses before it allows it to activate, which regulates the flow of the other
neurons in the stepped cluster

    19.158667ms [cycle] 19.158667ms
    35.780084ms [refraction] 35.780084ms
    2.169033709s [response] 542ns
    2.166928917s [cycle] 2.166928917s
    2.166973875s [refraction] 2.166967042s
    2.167370333s [response] 583ns
    2.166713708s [cycle] 2.166713708s
    2.166828958s [refraction] 2.166810375s
    2.166855042s [response] 750ns
    2.166907208s [cycle] 2.166907208s
    2.167134125s [refraction] 2.167114833s

I think that's a wonderful initial primer to the basics of using a neural cortex.  We'll start hitting some more
advanced topics in the next solution =)