# go-macos/sessionwatch

[![ci](https://github.com/go-macos/sessionwatch/actions/workflows/ci.yml/badge.svg)](https://github.com/go-macos/sessionwatch/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-macos/sessionwatch.svg)](https://pkg.go.dev/github.com/go-macos/sessionwatch)
[![License](https://img.shields.io/badge/license-BSD--3--Clause-blue.svg)](LICENSE)

**Know when the person at a Mac has stepped away, from pure Go
(`CGO_ENABLED=0`): the screen locked, the machine is going to sleep, or the
login session was switched out.** A program that holds something unlocked (a
decrypted container, an open vault, a mounted volume) closes it on the first
event. Extracted from [theca](https://github.com/go-theca/theca)'s
close-on-lock.

```go
ch, err := sessionwatch.Watch(ctx)
if err != nil {
	return err // sessionwatch.ErrUnsupported off macOS
}
go func() {
	if ev, ok := <-ch; ok { // sessionwatch.Lock, Sleep or Session
		closeEverything(ev)
	}
}()
objc.RunMainLoop(ctx) // the notifications are delivered on the main thread
```

## API

| | |
|---|---|
| `Watch(ctx) (<-chan Event, error)` | One `Event`, then the channel closes; it also closes, without an event, when `ctx` ends. |
| `Lock`, `Sleep`, `Session` | `com.apple.screenIsLocked`; `NSWorkspaceWillSleepNotification`; `NSWorkspaceSessionDidResignActiveNotification`. |
| `ErrUnsupported` | Off darwin, `Watch` returns it. |
| `Logf`, `LogTo(w)` | Say which mode is active (silent by default). |

## How it hears

Notifications go through [`go-macos/objc`](https://github.com/go-macos/objc)
`Observe`. They arrive before the sleep and without polling lag, but only while
the main thread runs its run loop (`objc.RunMainLoop`). `Watch` proves the
observers by posting a notification of its own and waiting for it to come back;
if registration fails, or nothing delivers (no run loop), it falls back to
polling `CGSSessionScreenIsLocked` every 2 s and noticing a sleep as a gap
between the wall clock and the monotonic clock (which stops during sleep on
macOS). `Logf` says which mode was chosen. A session without the CoreGraphics
dictionary (ssh) reads as not locked.

## Every platform

Every symbol exists everywhere; off darwin `Watch` returns `ErrUnsupported`,
so a consumer cross-compiles without build tags. Linux (logind over D-Bus) and
Windows (session notifications) are not covered here.

## Tests

Coverage is **100 %**, statement for statement, on darwin and off it. The
notification source is a seam the tests replace; one test registers the real
observers on the real main-thread run loop and proves them with the package's
own probe. No test posts or simulates a real lock or sleep.

## Install

```
go get github.com/go-macos/sessionwatch
```

BSD-3-Clause.
