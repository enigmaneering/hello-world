package main

import (
	"time"

	"git.enigmaneering.net/hello-world/enigma0/solution1/std"
	"git.ignitelabs.net/janos/core"
	"git.ignitelabs.net/janos/core/sys/atlas"
	"git.ignitelabs.net/janos/core/sys/rec"
)

func init() {
	core.Describe("contrived demonstration")
	atlas.Verbose(true)
}

func main() {
	counter := std.NewSynapse("Counter", func(imp *std.Impulse) {
		thought := imp.Manifest(0)

		count := thought.Revelation.(int)
		rec.Printf(imp.String(), "%v\n", count)
		thought.Revelation = count + 1
	}, nil, func(imp *std.Impulse) {
		go core.ShutdownNow()
	})

	toggler := std.NewSynapse("Toggler", func(imp *std.Impulse) {
		thought := imp.Manifest(true)

		toggle := thought.Revelation.(bool)
		if toggle {
			rec.Printf(imp.String(), "muting\n")
			counter <- std.Signal.Mute()
		} else {
			rec.Printf(imp.String(), "unmuting\n")
			counter <- std.Signal.Unmute()
		}
		thought.Revelation = !toggle
	}, func(imp *std.Impulse) bool {
		if imp.Beat == 3 {
			imp.Beat = 0
			return true
		}
		return false
	})

	counter <- std.Signal.Decay(7)

	source := std.Signal.Spark(&std.Impulse{
		Entity: ,
	})

	go func() {
		for core.Alive() {
			counter <- std.Signal.Spark(source)
			toggler <- std.Signal.Spark(source)
			time.Sleep(time.Second)
		}
	}()

	core.KeepAlive()
}
