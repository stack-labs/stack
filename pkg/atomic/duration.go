package atomic

import (
	"time"
)

// Duration is an atomic type-safe wrapper for time.Duration values.
type Duration struct {
	v Int64
}

var _zeroDuration time.Duration

// NewDuration creates a new Duration.
func NewDuration(v time.Duration) *Duration {
	x := &Duration{}
	if v != _zeroDuration {
		x.Store(v)
	}
	return x
}

// Load atomically loads the wrapped time.Duration.
func (d *Duration) Load() time.Duration {
	return time.Duration(d.v.Load())
}

// Store atomically stores the passed time.Duration.
func (d *Duration) Store(v time.Duration) {
	d.v.Store(int64(v))
}

// CAS is an atomic compare-and-swap for time.Duration values.
func (d *Duration) CAS(o, n time.Duration) bool {
	return d.v.CAS(int64(o), int64(n))
}

// Swap atomically stores the given time.Duration and returns the old
// value.
func (d *Duration) Swap(o time.Duration) time.Duration {
	return time.Duration(d.v.Swap(int64(o)))
}

// Add atomically adds to the wrapped time.Duration and returns the new value.
func (d *Duration) Add(n time.Duration) time.Duration {
	return time.Duration(d.v.Add(int64(n)))
}

// Sub atomically subtracts from the wrapped time.Duration and returns the new value.
func (d *Duration) Sub(n time.Duration) time.Duration {
	return time.Duration(d.v.Sub(int64(n)))
}

// String encodes the wrapped value as a string.
func (d *Duration) String() string {
	return d.Load().String()
}
