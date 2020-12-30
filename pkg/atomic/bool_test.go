package atomic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBool(t *testing.T) {
	atom := NewBool(false)
	require.False(t, atom.Toggle(), "Expected Toggle to return previous value.")
	require.True(t, atom.Toggle(), "Expected Toggle to return previous value.")
	require.False(t, atom.Toggle(), "Expected Toggle to return previous value.")
	require.True(t, atom.Load(), "Unexpected state after swap.")

	require.True(t, atom.CAS(true, true), "CAS should swap when old matches")
	require.True(t, atom.Load(), "CAS should have no effect")
	require.True(t, atom.CAS(true, false), "CAS should swap when old matches")
	require.False(t, atom.Load(), "CAS should have modified the value")
	require.False(t, atom.CAS(true, false), "CAS should fail on old mismatch")
	require.False(t, atom.Load(), "CAS should not have modified the value")

	atom.Store(false)
	require.False(t, atom.Load(), "Unexpected state after store.")

	prev := atom.Swap(false)
	require.False(t, prev, "Expected Swap to return previous value.")

	prev = atom.Swap(true)
	require.False(t, prev, "Expected Swap to return previous value.")

	t.Run("String", func(t *testing.T) {
		t.Run("true", func(t *testing.T) {
			assert.Equal(t, "true", NewBool(true).String(),
				"String() returned an unexpected value.")
		})

		t.Run("false", func(t *testing.T) {
			var b Bool
			assert.Equal(t, "false", b.String(),
				"String() returned an unexpected value.")
		})
	})
}
