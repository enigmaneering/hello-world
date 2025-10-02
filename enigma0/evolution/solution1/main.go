package main

import (
	"time"

	"git.enigmaneering.net/hello-world/enigma0/solution1/std"
	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/sys/rec"
)

func main() {
	syn := std.NewSynapse("Counter", func(imp *std.Impulse) {
		if imp.Thought == nil {
			imp.Thought = std.NewThought(0)
		}

		count := imp.Thought.Revelation.(int)
		rec.Printf(imp.String(), "%v\n", count)
		imp.Thought.Revelation = count + 1
	}, nil, func(imp *std.Impulse) {
		rec.Printf(imp.String(), "decayed\n")
	})

	syn <- std.Signal.Decay(4)

	for core.Alive() {
		syn <- std.Signal.Spark()
		time.Sleep(time.Second)
	}
}
