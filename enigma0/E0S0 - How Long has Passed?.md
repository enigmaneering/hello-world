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
this case we can simply rely upon the _phase_ of the cortical beat.  On every impulse, the cortex increments a beat
in the closed interval of `[0, Phase]` (inclusive).  If you set the phase to a negative value, the beat will infinitely
increment

    func main() {
        c := std.NewCortex(std.RandomName())
        c.Frequency = 60 //hz
        c.Phase = 3
    
        c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Cycle Printer", PrintCycle, func(imp *std.Impulse) bool {
            return imp.Beat == 0
        })
        c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Refractory Printer", PrintRefractory, func(imp *std.Impulse) bool {
            return imp.Beat == 1
        })
    
        c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Response Time Printer", PrintResponse, func(imp *std.Impulse) bool {
            return imp.Beat == 2
        })
    
        c.Spark()
        core.KeepAlive()
    }
    
    func PrintCycle(imp *std.Impulse) {
        fmt.Printf("%v [cycle] %v\n", imp.Timeline.CyclePeriod().String(), imp.Timeline.CyclePeriod().String())
    }
    
    func PrintRefractory(imp *std.Impulse) {
        fmt.Printf("%v [refraction] %v\n", imp.Timeline.CyclePeriod().String(), imp.Timeline.RefractoryPeriod().String())
    }
    
    func PrintResponse(imp *std.Impulse) {
        fmt.Printf("%v [response] %v\n", imp.Timeline.CyclePeriod().String(), imp.Timeline.ResponseTime().String())
    }


I've added the cycle period to the left so you can see the frequency of each function's activation: 20hz.  That's because
we're phasing at ⅓ the rate of 60hz.  You'll also notice the initial activation happens as fast as possible before the
remaining activations stabilize

    19.35025ms [cycle] 19.35025ms
    35.940209ms [refraction] 35.940209ms
    51.671459ms [response] 167ns
    66.620417ms [cycle] 66.620417ms
    66.700791ms [refraction] 66.678583ms
    66.839375ms [response] 42ns
    66.662583ms [cycle] 66.662583ms
    66.659792ms [refraction] 66.640792ms

The potential function, however, can be anything you'd desire - let's put a little faux delay in the last activation and
then bring the rate of speed down to 2hz

    c.Frequency = 2 //hz

    ...

	c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Cycle Printer", PrintCycle, makePotential(0))
	c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Refractory Printer", PrintRefractory, makePotential(1))

	i := 0
	c.Synapses() <- std.NewSynapse(lifecycle.Looping, "Response Time Printer", PrintResponse, func(imp *std.Impulse) bool {
		if imp.Beat == 2 {
			i++
			if i >= 3 {
				i = 0
				return true
			}
		}
		return false
	})

Here the final potential requires 3 impulses before it activates, which you can visibly watch time out in your console
at this slower speed

    503.240083ms [cycle] 503.240083ms
    1.003172583s [refraction] 1.003172583s
    1.999216s [cycle] 1.999216s
    2.000086084s [refraction] 2.00003925s
    2.00078625s [cycle] 2.00078625s
    2.000022791s [refraction] 1.99997875s
    5.502637s [response] 1.083µs

I think that's a wonderful initial primer to the basics of using a neural cortex.  We'll start hitting some more
advanced topics in the next solution =)