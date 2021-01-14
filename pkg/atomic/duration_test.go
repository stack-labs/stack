package atomic

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDuration(t *testing.T) {
	atom := NewDuration(5 * time.Minute)

	require.Equal(t, 5*time.Minute, atom.Load(), "Load didn't work.")
	require.Equal(t, 6*time.Minute, atom.Add(time.Minute), "Add didn't work.")
	require.Equal(t, 4*time.Minute, atom.Sub(2*time.Minute), "Sub didn't work.")

	require.True(t, atom.CAS(4*time.Minute, time.Minute), "CAS didn't report a swap.")
	require.Equal(t, time.Minute, atom.Load(), "CAS didn't set the correct value.")

	require.Equal(t, time.Minute, atom.Swap(2*time.Minute), "Swap didn't return the old value.")
	require.Equal(t, 2*time.Minute, atom.Load(), "Swap didn't set the correct value.")

	atom.Store(10 * time.Minute)
	require.Equal(t, 10*time.Minute, atom.Load(), "Store didn't set the correct value.")

	t.Run("String", func(t *testing.T) {
		assert.Equal(t, "42s", NewDuration(42*time.Second).String(),
			"String() returned an unexpected value.")
	})
}
