package darwin

import "C"
import (
	"sync"
	"unsafe"

	"git.ignitelabs.net/janos/core/std"
)

type Window struct {
	id        int64
	handle    unsafe.Pointer
	width     uint
	height    uint
	focused   bool
	minimized bool
	closed    bool
	x         uint
	y         uint

	mutex sync.Mutex

	Position std.TemporalBuffer[]
}

func (w *Window) cleanup() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.closed {
		return
	}
	w.closed = true

	winMux.Lock()
	delete(windows, w.id)
	windowCount := len(windows)
	winMux.Unlock()

	w.on.Broadcast()
	w.resize.Broadcast()
	w.move.Broadcast()
	w.focus.Broadcast()
	w.close.Broadcast()
	w.minimize.Broadcast()

	if windowCount == 0 {
		select {
		case allClosed <- struct{}{}:
		default:
		}
	}
}

//func (w *Window) Close() {
//	if w.closed {
//		return
//	}
//
//	// Destroy the NSWindow
//	C.destroyWindow(w.handle)
//
//	// Run cleanup
//	w.cleanup()
//}
