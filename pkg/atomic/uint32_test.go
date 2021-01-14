package atomic

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUint32(t *testing.T) {
	atom := NewUint32(42)

	require.Equal(t, uint32(42), atom.Load(), "Load didn't work.")
	require.Equal(t, uint32(46), atom.Add(4), "Add didn't work.")
	require.Equal(t, uint32(44), atom.Sub(2), "Sub didn't work.")
	require.Equal(t, uint32(45), atom.Inc(), "Inc didn't work.")
	require.Equal(t, uint32(44), atom.Dec(), "Dec didn't work.")

	require.True(t, atom.CAS(44, 0), "CAS didn't report a swap.")
	require.Equal(t, uint32(0), atom.Load(), "CAS didn't set the correct value.")

	require.Equal(t, uint32(0), atom.Swap(1), "Swap didn't return the old value.")
	require.Equal(t, uint32(1), atom.Load(), "Swap didn't set the correct value.")

	atom.Store(42)
	require.Equal(t, uint32(42), atom.Load(), "Store didn't set the correct value.")

	t.Run("String", func(t *testing.T) {
		// Use an integer with the signed bit set. If we're converting
		// incorrectly, we'll get a negative value here.
		atom := NewUint32(math.MaxUint32)
		assert.Equal(t, "4294967295", atom.String(),
			"String() returned an unexpected value.")
	})
}
