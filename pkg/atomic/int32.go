package atomic

import (
	"strconv"
	"sync/atomic"
)

// Int32 is an atomic wrapper around int32.
type Int32 struct {
	v int32
}

// NewInt32 creates a new Int32.
func NewInt32(i int32) *Int32 {
	return &Int32{v: i}
}

// Load atomically loads the wrapped value.
func (i *Int32) Load() int32 {
	return atomic.LoadInt32(&i.v)
}

// Add atomically adds to the wrapped int32 and returns the new value.
func (i *Int32) Add(n int32) int32 {
	return atomic.AddInt32(&i.v, n)
}

// Sub atomically subtracts from the wrapped int32 and returns the new value.
func (i *Int32) Sub(n int32) int32 {
	return atomic.AddInt32(&i.v, -n)
}

// Inc atomically increments the wrapped int32 and returns the new value.
func (i *Int32) Inc() int32 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped int32 and returns the new value.
func (i *Int32) Dec() int32 {
	return i.Sub(1)
}

// CAS is an atomic compare-and-swap.
func (i *Int32) CAS(old, new int32) bool {
	return atomic.CompareAndSwapInt32(&i.v, old, new)
}

// Store atomically stores the passed value.
func (i *Int32) Store(n int32) {
	atomic.StoreInt32(&i.v, n)
}

// Swap atomically swaps the wrapped int32 and returns the old value.
func (i *Int32) Swap(n int32) int32 {
	return atomic.SwapInt32(&i.v, n)
}

// String encodes the wrapped value as a string.
func (i *Int32) String() string {
	v := i.Load()
	return strconv.FormatInt(int64(v), 10)
}
