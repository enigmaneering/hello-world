# `E0S3 - Thoughts`
### `Alex Petz, Ignite Laboratories, September 2025`

---

### Neural Servers

Let's pose ourselves a challenge

    Set up an environment that can self-launch, shutdown, and restart multiple web servers at once

This sounds simple, on its face - self-restarting web servers!  That's nothing new!  And certainly, you can
achieve this with the existing tools.  Let's look at how a neural architecture might approach this, though.

This solution is a little more deep than those before it, but I'll walk you through each point step-by-step.  First, 
let's discuss _what's_ happening in this example.  At an even interval the cortex launches off _three_ web server 
synapses, but only _one_ of them will succeed in grabbing port _4242_ at a time.  The two that don't grab the port will 
cycle out and wait for the next available impulse. On activation, the impulse forms a _thought_ - quite literally!  This 
thought is a reference held by all future impulses, allowing future activations to mature it over time.  All thoughts, in
addition, are paired with a 'gate' mutex for thread-safe access.

    // Action Function

    // 0 - Create the server

    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        _, _ = w.Write([]byte(fmt.Sprintf("\"Hello, World!\"\n\t- %v", imp.Bridge)))
    })
    server := &http.Server{
        Addr:    ":" + strconv.Itoa(port),
        Handler: handler,
    }

    // 1 - Assign it to the impulse's 'thought'

    imp.Thought = std.NewThought(server)

From there, we fire off two asynchronous operations.  First, we launch the thought's web server in a goroutine

    go func() {
        rec.Printf(imp.Bridge, "launching server at :%d\n", port)

        if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
            rec.Printf(imp.Bridge, "disconnected\n") // This is an intended 'shutdown' event
        } else {
            rec.Printf(imp.Bridge, "cycling\n") // This is an error
        }

        // Release the thought
        imp.Thought = nil
    }()

When the server disconnects, it sets the thought to nil - allowing it to be garbage collected and signal for other
synapses to begin activating.  Next, we fire off a faux delayed shutdown of the web server to demonstrate our cyclic
design

    go func() {
        time.Sleep(time.Second * 2)
        if imp.Thought != nil {
            // Access to thoughts can be gated for thread safety
            imp.Thought.Gate.Lock()
            defer imp.Thought.Gate.Unlock()

            _ = imp.Thought.Revelation.(*http.Server).Shutdown(context.Background())
        }
    }()

Moving down from there we reach our potential function, which simply stops future impulses when the server is already running
on that neural endpoint

    func Potential(imp *std.Impulse) bool {
        // Only serve when the impulsive thought isn't "busy"
        if imp.Thought == nil {
            return true
        }
        return false
    }

Lastly, we reach the neuron's cleanup function.  This is what allows system shutdown events to properly clean up the 
web server's resources

    func Cleanup(imp *std.Impulse) {
        // Access the thought and shut it down
        if imp.Thought != nil {
            imp.Thought.Gate.Lock()
            defer imp.Thought.Gate.Unlock()

            _ = imp.Thought.Revelation.(*http.Server).Shutdown(context.Background())
        }
    }

This is our resulting output

    [Smedley Kynforth ⇝ Server C] launching server at :4242
    [Smedley Kynforth ⇝ Server A] launching server at :4242
    [Smedley Kynforth ⇝ Server B] launching server at :4242
    [Smedley Kynforth ⇝ Server A] cycling
    [Smedley Kynforth ⇝ Server B] cycling
    [Smedley Kynforth ⇝ Server A] launching server at :4242
    [Smedley Kynforth ⇝ Server B] launching server at :4242
    [Smedley Kynforth ⇝ Server B] cycling
    [Smedley Kynforth ⇝ Server A] cycling
    [Smedley Kynforth ⇝ Server C] disconnected
    [Smedley Kynforth ⇝ Server B] launching server at :4242
    [Smedley Kynforth ⇝ Server A] launching server at :4242
    [Smedley Kynforth ⇝ Server B] cycling
    [Smedley Kynforth ⇝ Server C] launching server at :4242
    [Smedley Kynforth ⇝ Server C] cycling

More importantly, if you ping `localhost:4242` in your web browser you'll be greeted with your synaptic
endpoint, identifying itself amongst the void of internet traffic - a lone sentinel of pleasantries and cordiality
in an ever-growing sea of neural cohorts, empathetically holding the line between exuberance and apathy

    "Hello, World!"
	    - Beag Hessenthaler ⇝ Server B

_**..._or no page at all! =)_**_

_(if pinged between activations)_

### Findings

We haven't met our challenge, yet!  This is a terribly inefficient design with _zero_ utility at the moment,
but by changing just two lines of code we resolve the efficiency _and_ solve the challenge!

	cortex.Synapses() <- std.NewSynapse(lifecycle.Looping, "Server A", Serve(4242), Potential, Cleanup)

                                                                              ⬐ new port 
	cortex.Synapses() <- std.NewSynapse(lifecycle.Looping, "Server B", Serve(4243), Potential, Cleanup)
	cortex.Synapses() <- std.NewSynapse(lifecycle.Looping, "Server C", Serve(4244), Potential, Cleanup)
                                                                     new port ⬏

Now, rather than fighting for a shared resource, each owns a specific network endpoint and cycles itself automatically!
This lets things run significantly more "stably" - albiet absurdly so.

This is a _very_ wordy way of doing something simple, though - and the entire point of aneural architecture is to
mitigate the amount of repeated code.  If you'd like to simply create a basic neural web server which can reference
a thought between requests, you'll find use in the `neural` package.  Rather than recreating the logic in this solution,
you can create the entire stack by defining your synapse as such

	cortex.Synapses() <- neural.Net.Server(std.RandomName(), ":4242", func(imp *std.Impulse) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            _, _ = w.Write([]byte(fmt.Sprintf("\"Hello, World!\"\n\t- %v", imp.Bridge)))
        })
	})

Nothing else required =)