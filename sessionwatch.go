// Package sessionwatch tells a process when the person at a Mac has stepped
// away: the screen locked, the machine is going to sleep, or the login session
// was switched out. A program that holds something unlocked (a decrypted
// container, an open vault, a mounted volume) closes it on the first event.
//
// It listens for com.apple.screenIsLocked (NSDistributedNotificationCenter),
// NSWorkspaceWillSleepNotification and
// NSWorkspaceSessionDidResignActiveNotification through go-macos/objc. They
// arrive before the sleep and without polling lag, on the main thread's run
// loop, which the caller runs (objc.RunMainLoop). Watch proves the observers
// with a notification of its own; when they cannot be registered or heard it
// falls back to reading the session's CGSSessionScreenIsLocked every 2 s and
// noticing a sleep as a gap between the wall clock and the monotonic clock.
//
// Every symbol exists on every platform; off darwin Watch returns
// ErrUnsupported, so a cross-platform consumer needs no build tags.
package sessionwatch

import (
	"context"
	"errors"
	"fmt"
	"io"
)

// Event is the reason the person has gone.
type Event string

// The events Watch delivers.
const (
	// Lock: the screen was locked.
	Lock Event = "lock"
	// Sleep: the machine is going to sleep.
	Sleep Event = "sleep"
	// Session: the login session was switched out (fast user switching).
	Session Event = "session"
)

// ErrUnsupported is returned where no lock or sleep source is implemented.
var ErrUnsupported = errors.New("sessionwatch: no screen-lock or sleep notification on this system")

// Logf says which detection mode is active (notifications or polling). It is
// silent by default.
var Logf = func(format string, args ...any) {}

// LogTo makes Logf write one line per message to w.
func LogTo(w io.Writer) {
	Logf = func(format string, args ...any) {
		fmt.Fprintf(w, format+"\n", args...)
	}
}

// Watch returns a channel that receives one Event when the person has gone and
// is then closed; it is also closed, without an event, when ctx ends. On macOS
// the notifications need the main thread's run loop (objc.RunMainLoop).
func Watch(ctx context.Context) (<-chan Event, error) { return watch(ctx) }
