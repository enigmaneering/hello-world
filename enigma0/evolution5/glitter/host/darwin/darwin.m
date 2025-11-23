#import <Cocoa/Cocoa.h>

void initApp() {
    [NSApplication sharedApplication];
    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
    [NSApp activateIgnoringOtherApps:YES];
}

// Forward declaration - now takes long instead of void*
void goWindowEvent(long windowID, int eventType, int data1, int data2);

// Window delegate that forwards to Go
@interface WindowDelegate : NSObject <NSWindowDelegate>
@property (assign) long windowID;
@end

@implementation WindowDelegate

- (void)windowDidResize:(NSNotification*)notification {
    NSWindow* window = [notification object];
    NSRect frame = [window frame];
    goWindowEvent(self.windowID, 1, (int)frame.size.width, (int)frame.size.height);
}

- (void)windowDidMove:(NSNotification*)notification {
    NSWindow* window = [notification object];
    NSRect frame = [window frame];
    goWindowEvent(self.windowID, 2, (int)frame.origin.x, (int)frame.origin.y);
}

- (void)windowDidBecomeKey:(NSNotification*)notification {
    goWindowEvent(self.windowID, 3, 0, 0);  // Focus gained
}

- (void)windowDidResignKey:(NSNotification*)notification {
    goWindowEvent(self.windowID, 4, 0, 0);  // Focus lost
}

- (BOOL)windowShouldClose:(NSWindow*)sender {
    goWindowEvent(self.windowID, 5, 0, 0);  // Close requested
    return YES;  // Allow close
}

- (void)windowWillClose:(NSNotification*)notification {
    goWindowEvent(self.windowID, 6, 0, 0);  // Window closing
}

- (void)windowDidMiniaturize:(NSNotification*)notification {
    goWindowEvent(self.windowID, 7, 0, 0);  // Minimized
}

- (void)windowDidDeminiaturize:(NSNotification*)notification {
    goWindowEvent(self.windowID, 8, 0, 0);  // Restored
}

@end

void* createWindow(int x, int y, int width, int height, const char* title, long windowID) {
    NSRect frame = NSMakeRect(x, y, width, height);

    NSWindowStyleMask style = NSWindowStyleMaskTitled |
                              NSWindowStyleMaskClosable |
                              NSWindowStyleMaskMiniaturizable |
                              NSWindowStyleMaskResizable;

    NSWindow* window = [[NSWindow alloc]
        initWithContentRect:frame
        styleMask:style
        backing:NSBackingStoreBuffered
        defer:NO];

    [window setTitle:[NSString stringWithUTF8String:title]];

    window.titlebarAppearsTransparent = YES;

    // Create and set delegate
    WindowDelegate* delegate = [[WindowDelegate alloc] init];
    delegate.windowID = windowID;
    [window setDelegate:delegate];

    [window makeKeyAndOrderFront:nil];

    return (void*)window;
}

// Simple event loop - just keeps app alive
void runEventLoop() {
    while (true) {
        @autoreleasepool {
            NSEvent* event = [NSApp nextEventMatchingMask:NSEventMaskAny
                                                untilDate:[NSDate distantFuture]
                                                   inMode:NSDefaultRunLoopMode
                                                  dequeue:YES];
            [NSApp sendEvent:event];
            [NSApp updateWindows];
        }
    }
}

void destroyWindow(void* windowPtr) {
    if (!windowPtr) return;

    NSWindow* window = (NSWindow*)windowPtr;

    if ([window isVisible]) {
        [window close];
    }

    WindowDelegate* delegate = (WindowDelegate*)[window delegate];
    [window setDelegate:nil];
}
