package std

import (
	"runtime"
	"sync"
	"time"
)

// A Gate can patiently Attempt to interface with a sync.Mutex.
type Gate struct {
	sync.Mutex
}

// Attempt repeatedly tries to "hold the gate" before conceding if unable to do so in the allotted time.
// This will return true if a lock was attained - otherwise, false.
//
// NOTE: This will attempt 111 locks (completely arbitrary) before dropping to "lock polling", which starts at
// 1µs and exponentially decays the polling rate of every cycle up to a maximum of 11ms.  Theoretically, this should
// cause just over two hundred lock attempts in the first second before settling in at roughly ninety a second.
func (g *Gate) Attempt(timeout time.Duration) bool {
	start := time.Now()
	for i := 0; i < 111; i++ {
		if g.TryLock() {
			return true
		}
		if i%11 == 0 {
			// Don't be greedy!
			runtime.Gosched()
		}
		if time.Since(start) >= timeout {
			return false
		}
	}

	delay := 1 * time.Microsecond
	maxDelay := 11 * time.Millisecond

	for {
		if time.Since(start) >= timeout {
			return false
		}

		if g.TryLock() {
			return true
		}

		time.Sleep(delay)
		if delay < maxDelay {
			delay *= 2
			if delay > maxDelay {
				delay = maxDelay
			}
		}
	}
}
