package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/enum/lifecycle"
	"git.ignitelabs.net/janos/core/std"
	"git.ignitelabs.net/janos/core/sys/rec"
)

/*
E0S4

This creates three self-recycling web servers.
*/

var cortex = std.NewCortex(std.RandomName())

func main() {
	cortex.Frequency = 4 //hz

	cortex.Synapses() <- std.NewSynapse(lifecycle.Looping, "Server A", Serve(4242), Potential, Cleanup)
	cortex.Synapses() <- std.NewSynapse(lifecycle.Looping, "Server B", Serve(4242), Potential, Cleanup)
	cortex.Synapses() <- std.NewSynapse(lifecycle.Looping, "Server C", Serve(4242), Potential, Cleanup)

	// NOTE: Set each synapse to different ports for a more stable cycle =)

	cortex.Spark()
	core.KeepAlive(time.Second * 2)
}

func Serve(port int) func(imp *std.Impulse) {
	return func(imp *std.Impulse) {
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

		// 2 - Launch the server asynchronously

		go func() {
			rec.Printf(imp.Bridge, "launching server at :%d\n", port)

			if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
				rec.Printf(imp.Bridge, "disconnected\n")
			} else {
				rec.Printf(imp.Bridge, "cycling\n")
			}

			// Release the thought
			imp.Thought = nil
		}()

		// 3 - Create a faux delayed shutdown

		go func() {
			time.Sleep(time.Second * 2)
			if imp.Thought != nil {
				// Access to thoughts can be gated for thread safety
				imp.Thought.Gate.Lock()
				defer imp.Thought.Gate.Unlock()

				_ = imp.Thought.Revelation.(*http.Server).Shutdown(context.Background())
			}
		}()
	}
}

func Potential(imp *std.Impulse) bool {
	// Only serve when the impulsive thought isn't "busy"
	if imp.Thought == nil {
		return true
	}
	return false
}

func Cleanup(imp *std.Impulse) {
	// Access the thought and shut it down
	if imp.Thought != nil {
		imp.Thought.Gate.Lock()
		defer imp.Thought.Gate.Unlock()

		_ = imp.Thought.Revelation.(*http.Server).Shutdown(context.Background())
	}
}
