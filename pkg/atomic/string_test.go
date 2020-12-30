package atomic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringNoInitialValue(t *testing.T) {
	atom := &String{}
	require.Equal(t, "", atom.Load(), "Initial value should be blank string")
}

func TestString(t *testing.T) {
	atom := NewString("")
	require.Equal(t, "", atom.Load(), "Expected Load to return initialized value")

	atom.Store("abc")
	require.Equal(t, "abc", atom.Load(), "Unexpected value after Store")

	atom = NewString("bcd")
	require.Equal(t, "bcd", atom.Load(), "Expected Load to return initialized value")

	atom = NewString("foo")

	t.Run("String", func(t *testing.T) {
		atom := NewString("foo")
		assert.Equal(t, "foo", atom.String(),
			"String() returned an unexpected value.")
	})
}
