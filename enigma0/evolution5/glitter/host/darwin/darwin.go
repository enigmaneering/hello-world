package darwin

/*
#cgo darwin CFLAGS: -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

void initApp();
void* createWindow(int x, int y, int width, int height, const char* title, long windowID);
void runEventLoop();
*/
import "C"
import (
	"runtime"
	"sync"
	"unsafe"

	"git.enigmaneering.net/hello-world/enigma0/solution0/evolution5/glitter/host/darwin/event"
)

var (
	windows   = make(map[int64]*Window)
	winMux    sync.RWMutex
	nextWinID int64 = 1
)

var allClosed = make(chan any, 1<<10)
var AllClosed = <-allClosed

//export goWindowEvent
func goWindowEvent(windowID C.long, eventType C.int, data1 C.int, data2 C.int) {
	winMux.RLock()
	win, exists := windows[int64(windowID)]
	winMux.RUnlock()

	if !exists {
		return
	}

	win.on.Lock()
	msg := event.Message(eventType)
	switch msg {
	case event.Resize:
		win.resize.Lock()
		win.width = uint(data1)
		win.height = uint(data2)
		win.on.msg = Message{
			Message: msg,
			Data1:   win.width,
			Data2:   win.height,
		}
		win.resize.Broadcast()
		win.resize.Unlock()

	case event.Move:
		win.move.Lock()
		win.x = uint(data1)
		win.y = uint(data2)
		win.on.msg = Message{
			Message: msg,
			Data1:   win.x,
			Data2:   win.y,
		}
		win.move.Broadcast()
		win.move.Unlock()

	case event.FocusGain:
		win.focus.Lock()
		win.focused = true
		win.on.msg = Message{
			Message: msg,
		}
		win.focus.Broadcast()
		win.focus.Unlock()

	case event.FocusLose:
		win.focus.Lock()
		win.focused = false
		win.on.msg = Message{
			Message: msg,
		}
		win.focus.Broadcast()
		win.focus.Unlock()

	case event.Closed, event.CloseReq:
		win.close.Lock()
		win.on.msg = Message{
			Message: msg,
		}
		win.close.Broadcast()
		win.close.Unlock()
		win.cleanup()

	case event.Minimize:
		win.minimize.Lock()
		win.minimized = true
		win.on.msg = Message{
			Message: msg,
		}
		win.minimize.Broadcast()
		win.minimize.Unlock()

	case event.Restore:
		win.minimize.Lock()
		win.minimized = false
		win.on.msg = Message{
			Message: msg,
		}
		win.minimize.Broadcast()
		win.minimize.Unlock()
	}
	win.on.Broadcast()
	win.on.Unlock()
}

func init() {
	runtime.LockOSThread()
}

func CreateWindow(x, y, width, height int, title string) *Window {
	winMux.Lock()
	windowID := nextWinID
	nextWinID++
	winMux.Unlock()

	win := &Window{
		id:       windowID,
		width:    uint(width),
		height:   uint(height),
		x:        uint(x),
		y:        uint(y),
		on:       newCondition(),
		resize:   newCondition(),
		move:     newCondition(),
		focus:    newCondition(),
		minimize: newCondition(),
		close:    newCondition(),
	}
	win.Delegate = delegator{win: any(win)}

	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))

	// Pass integer ID instead of Go pointer
	handle := C.createWindow(C.int(x), C.int(y), C.int(width), C.int(height), cTitle, C.long(windowID))

	win.handle = handle

	winMux.Lock()
	windows[windowID] = win
	winMux.Unlock()

	return win
}

func Start() {
	C.initApp()
	C.runEventLoop()
}
