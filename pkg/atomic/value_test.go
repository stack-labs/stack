package atomic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	var v Value
	assert.Nil(t, v.Load(), "initial Value is not nil")

	v.Store(42)
	assert.Equal(t, 42, v.Load())

	v.Store(84)
	assert.Equal(t, 84, v.Load())

	assert.Panics(t, func() { v.Store("foo") })
}
