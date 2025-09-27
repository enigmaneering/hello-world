package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/std"
	"git.ignitelabs.net/janos/core/std/neural"
	"git.ignitelabs.net/janos/core/sys/rec"
)

/*
E0S4 - Simple

This demonstrates a simple self-restarting neural.Net.Server.
*/

func main() {
	cortex := std.NewCortex(std.RandomName())
	cortex.Frequency = 1 //hz
	cortex.Mute()

	port := ":4242"

	// NOTE: Change this to `Handler` for a simple stable server demonstration ⬎
	cortex.Synapses() <- neural.Net.Server("localhost"+port, port, HandlerWhichShutsDown, func(imp *std.Impulse) {
		cortex.Impulse()
	})

	cortex.Spark()
	cortex.Impulse()
	core.KeepAlive(time.Second)
}

func Handler(imp *std.Impulse) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf("\"Hello, World!\"\n\t- %v", imp.Bridge)))
	})
}

func HandlerWhichShutsDown(imp *std.Impulse) http.Handler {
	// Introduce a faux delayed shutdown
	go func() {
		delay := time.Second * 5
		rec.Printf(imp.Bridge, "disconnecting in %v\n", delay)
		time.Sleep(delay)
		if imp.Thought != nil && imp.Thought.Revelation != nil {
			_ = imp.Thought.Revelation.(*http.Server).Shutdown(context.Background())
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(fmt.Sprintf("\"Hello, World!\"\n\t- %v", imp.Bridge)))
	})
}
