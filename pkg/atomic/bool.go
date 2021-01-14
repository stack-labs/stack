package atomic

import "strconv"

// Bool is an atomic type-safe wrapper for bool values.
type Bool struct {
	v Uint32
}

var _zeroBool bool

// NewBool creates a new Bool.
func NewBool(v bool) *Bool {
	x := &Bool{}
	if v != _zeroBool {
		x.Store(v)
	}
	return x
}

// Load atomically loads the wrapped bool.
func (b *Bool) Load() bool {
	return truthy(b.v.Load())
}

// Store atomically stores the passed bool.
func (b *Bool) Store(v bool) {
	b.v.Store(boolToInt(v))
}

// CAS is an atomic compare-and-swap for bool values.
func (b *Bool) CAS(o, n bool) bool {
	return b.v.CAS(boolToInt(o), boolToInt(n))
}

// Swap atomically stores the given bool and returns the old
// value.
func (b *Bool) Swap(o bool) bool {
	return truthy(b.v.Swap(boolToInt(o)))
}

func truthy(n uint32) bool {
	return n == 1
}

func boolToInt(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

// Toggle atomically negates the Boolean and returns the previous value.
func (b *Bool) Toggle() bool {
	for {
		old := b.Load()
		if b.CAS(old, !old) {
			return old
		}
	}
}

// String encodes the wrapped value as a string.
func (b *Bool) String() string {
	return strconv.FormatBool(b.Load())
}
