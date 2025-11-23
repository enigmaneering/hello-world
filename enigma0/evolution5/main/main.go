package main

import (
	"fmt"

	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution5/glitter/host/darwin"
	"git.ignitelabs.net/janos/core"
)

func main() {
	fmt.Println("here")
	win := darwin.CreateWindow(100, 100, 640, 480, "My Window")
	darwin.CreateWindow(100, 100, 640, 480, "My Window2")

	win.Delegate.OnResize(func(width, height uint) {
		fmt.Printf("Window resize: %dx%d\n", width, height)
	})
	win.Delegate.OnMove(func(x, y uint) {
		fmt.Printf("Window move: %dx%d\n", x, y)
	})
	win.Delegate.OnFocus(func(focused bool) {
		fmt.Printf("Window focused: %v\n", focused)
	})
	win.Delegate.OnMinimize(func(minimized bool) {
		fmt.Printf("Window minimized: %v\n", minimized)
	})

	go func() {
		for core.Alive() {
			select {
			case msg := <-win.Delegate.On():
				fmt.Println(msg)
			}
		}
	}()

	//go func() {
	//	time.Sleep(time.Second * 5)
	//	win.Close()
	//}()

	darwin.Start()
}
