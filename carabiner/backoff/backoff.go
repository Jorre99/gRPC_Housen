// Package backoff implements backoff waiting logic.
package backoff

import "time"

// Backoff represents state of an exponential backoff system.
type Backoff struct {
	start time.Duration
	mul   time.Duration
	max   time.Duration
	curr  time.Duration
}

// NewBackoff create a new backoff struct.
func NewBackoff(start, max time.Duration, mul int) *Backoff {
	return &Backoff{
		start: start,
		max:   max,
		curr:  start,
		mul:   time.Duration(mul),
	}
}

// Incr increases the backoff time, this should be called after a failure.
func (b *Backoff) Incr() {
	b.curr = min(b.curr*b.mul, b.max)
}

// Reset resets the backoff time to the initial value, this should be called after success.
func (b *Backoff) Reset() {
	b.curr = b.start
}

// Get returns the required wait time until the next attempt.
func (b *Backoff) Get() time.Duration {
	return b.curr
}

func min(i, j time.Duration) time.Duration {
	if i < j {
		return i
	}
	return j
}
