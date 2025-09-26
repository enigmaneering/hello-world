# `E0 - The Neural Impulse Engine`
## `A.K.A. The Self-Regulating Looping Clock`
### `Alex Petz, Ignite Laboratories, March 2025`

---

### Jumping Gee Willikers - Looping Clocks!?

Well, yes!  And there's a fundamentally important reason for this:

_**Calculation takes time**_

Even the act of _reading_ a value takes time!  Subsequent layers of calculation only exacerbate
the issue, but intelligent systems _demand_ rich calculations.  At any given instant billions of neurons are
firing across your biological machinery - accumulating thresholds, locking semaphores, and calculating _thought._

_How!?_

Easy - by harnessing the _shared_ passing of time through _feedback loops._

### The Colonel of Kernels

The entirety of this project is built around a distributed operating system I've named `JanOS` - after the Roman
God of time, duality, and beginnings.  Every Go program that imports the library becomes a JanOS **instance** by extension,
regardless of its integration into a larger system.  Your instance is given a random name at startup, though you can
set the defaults through the `atlas` configuration system (which we'll discuss later).  The instance provides global 
synchronization over any neural activity you'd like to **spark** off. 

**Neural** coordination begins at the **cortex**.  A cortex acts as a pulse of time which can coordinate clusters of neural activity.
If we liken this to hardware terms, a cortex is an oscillating crystal - albeit with _variable_ frequency.  This frequency drives
the execution of neural activity, and JanOS provides it as a first-class citizen.  A cortex drives many **synapses**, which act
as axonal bridges between a **neuron** and a cortex (though many neurons can be impulsed by many cortices).  On every beat of the 
cortical clock, registered synapses create **impulse** objects which build context as they reach their neural endpoints.  Finally, 
an **action-potential** is reached by triggering a potential function before deciding to activate the associated action.

    tl;dr - JanOS is a neural goroutine coordinator

The most striking quality of cortices is that they can intelligently be **muted** and **unmuted**.  This allows you to wait for
a condition (such as a button press) _before_ stimulating neural activity.  A synapse is simply a function (leaving the door open
to create your own) which typically ends when its creating cortex shuts down.  It acts as a way to efficiently gate impulses
against the cortical clock, allowing the neural endpoint to focus entirely on the job at hand.  Out of the box, there are four
different synaptic **lifecycles** which you can create -

**0.** Looping - where the neuron **indefinitely** fires in an efficient cyclic loop while the potential is high.

**1.** Stimulative - where the neuron **indefinitely** fires for _every_ impulse the potential is high.

**2.** Triggered - where the neuron is fired exactly once after the potential goes high and then **decays**.

**3.** Impulsed - where the neuron is attemped to be fired exactly once regardless of potential and then **decays**.

Let's step through the neural lifecycles, starting with an impulse from the cortex

<picture>
<img alt="Single Beat" src="assets/E0S1D0%20-%20Single%20Beat.svg" width="250" style="display: block; margin-left: auto; margin-right: auto;">
</picture>

Each synapse tracks the impulses as 'beats' of the cortical clock - with every subsequent beat forming a loop

<picture>
<img alt="Multiple Beats" src="assets/E0S1D1 - Multiple Beats.svg" width="350" style="display: block; margin-left: auto; margin-right: auto;">
</picture>

For every impulse, a potential function is called.  If it returns true, the provided action is called.  Is it an even beat?  Odd?  Could it "skip"
every 15 beats by modulo-ing the current beat with 16?  These are some of the more primitive things a potential can do to pace itself across time.
Once the potential returns "high" (meaning true) the neural activation occurs.

<picture>
<img alt="Activation" src="assets/E0S1D2 - Activation.svg" width="400" style="display: block; margin-left: auto; margin-right: auto;">
</picture>

Here we can see a timeline describing an in-flight activation on beat 42.  The key thing you should notice is that a neural timing engine can only get your
activation _close enough_ to the target moment - as entropy (often in the form of the task scheduler) is unavoidable.  This is a _good thing,_ but
one that should be duly noted.  The impulse's **timeline** records the most important moments during the neural lifecycle, allowing you to observe
the **refractory period** of a neural loop.

<picture>
<img alt="Looping Activation" src="assets/E0S1D3 - Looping Activation.svg" width="500" style="display: block; margin-left: auto; margin-right: auto;">
</picture>

When using a _looping_ activation, the refractory period represents a proportionally diminishing value as calculation time increases (relative to the
impulse frequency) - allowing you to quickly see if your calculations have or will exceed that pace.  In our above example, the neural activation 
appears to calculate the beat number it was activated on - but the rate of impulse exceeds the calculation frequency.  To resolve that issue, you could
slow down the rate of impulse or 'block' the cortex by muting and unmuting it while your neuron activates.  This pattern is called 'blocking activation'

<picture>
<img alt="Blocking Activation" src="assets/E0S1D5 - Blocking Activation.svg" width="500" style="display: block; margin-left: auto; margin-right: auto;">
</picture>

A looping activation allows you to cache data on the impulse object as a **thought** - allowing you to perform differential analysis _on impulse._  If
you don't need those qualities, you could drop from a looping activation into a _stimulative_ activation

<picture>
<img alt="Stimulative Activation" src="assets/E0S1D4 - Stimulative Activation.svg" width="500" style="display: block; margin-left: auto; margin-right: auto;">
</picture>

Finally, the ultimate goal is to layer many tracks of neural execution on top of one another - allowing you to build _clustered_ activations

<picture>
<img alt="Clustered Activation" src="assets/E0S1D7 - Clustered Activation.svg" width="550" style="display: block; margin-left: auto; margin-right: auto;">
</picture>

When combining all of the above together, you can achieve what I call _temporal fragment shading_ through stimulated neural execution

<picture>
<img alt="Temporal Fragment Shading" src="assets/E0S1D6 - Logical Activation.svg" width="500" style="display: block; margin-left: auto; margin-right: auto;">
</picture>

But, I'm getting quite ahead of myself!  I just want you to be familiar with the concept of visualizing neural tracks of execution across time.  This
will _immensely_ help you mentally work through the concepts I'm about to start manifesting.