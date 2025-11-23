package glitter2

//
///*
//#cgo darwin CFLAGS: -x objective-c
//#cgo darwin LDFLAGS: -framework Cocoa
//
//#import <Cocoa/Cocoa.h>
//
//// Must run on main thread for Cocoa
//void* createWindow(int x, int y, int width, int height, const char* title) {
//    NSRect frame = NSMakeRect(x, y, width, height);
//
//    NSWindowStyleMask style = NSWindowStyleMaskTitled |
//                              NSWindowStyleMaskClosable |
//                              NSWindowStyleMaskMiniaturizable |
//                              NSWindowStyleMaskResizable;
//
//    NSWindow* window = [[NSWindow alloc]
//        initWithContentRect:frame
//        styleMask:style
//        backing:NSBackingStoreBuffered
//        defer:NO];
//
//    [window setTitle:[NSString stringWithUTF8String:title]];
//    [window makeKeyAndOrderFront:nil];
//
//    return (void*)window;
//}
//
//void destroyWindow(void* window) {
//    NSWindow* w = (NSWindow*)window;
//    [w close];
//}
//
//void setWindowPosition(void* window, int x, int y) {
//    NSWindow* w = (NSWindow*)window;
//    NSPoint point = NSMakePoint(x, y);
//    [w setFrameOrigin:point];
//}
//
//void setWindowSize(void* window, int width, int height) {
//    NSWindow* w = (NSWindow*)window;
//    NSRect frame = [w frame];
//    frame.size = NSMakeSize(width, height);
//    [w setFrame:frame display:YES];
//}
//*/
//import "C"
//import (
//	"unsafe"
//)
//
//type Window struct {
//	handle unsafe.Pointer
//	width  int
//	height int
//	x      int
//	y      int
//}
//
//func CreateWindow(x, y, width, height int, title string) *Window {
//	cTitle := C.CString(title)
//	defer C.free(unsafe.Pointer(cTitle))
//
//	handle := C.createWindow(C.int(x), C.int(y), C.int(width), C.int(height), cTitle)
//
//	return &Window{
//		handle: handle,
//		width:  width,
//		height: height,
//		x:      x,
//		y:      y,
//	}
//}
//
//func (w *Window) Close() {
//	C.destroyWindow(w.handle)
//}
//
//func (w *Window) SetPosition(x, y int) {
//	w.x = x
//	w.y = y
//	C.setWindowPosition(w.handle, C.int(x), C.int(y))
//}
//
//func (w *Window) SetSize(width, height int) {
//	w.width = width
//	w.height = height
//	C.setWindowSize(w.handle, C.int(width), C.int(height))
//}
