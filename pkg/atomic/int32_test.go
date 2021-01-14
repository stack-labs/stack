package atomic

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt32(t *testing.T) {
	atom := NewInt32(42)

	require.Equal(t, int32(42), atom.Load(), "Load didn't work.")
	require.Equal(t, int32(46), atom.Add(4), "Add didn't work.")
	require.Equal(t, int32(44), atom.Sub(2), "Sub didn't work.")
	require.Equal(t, int32(45), atom.Inc(), "Inc didn't work.")
	require.Equal(t, int32(44), atom.Dec(), "Dec didn't work.")

	require.True(t, atom.CAS(44, 0), "CAS didn't report a swap.")
	require.Equal(t, int32(0), atom.Load(), "CAS didn't set the correct value.")

	require.Equal(t, int32(0), atom.Swap(1), "Swap didn't return the old value.")
	require.Equal(t, int32(1), atom.Load(), "Swap didn't set the correct value.")

	atom.Store(42)
	require.Equal(t, int32(42), atom.Load(), "Store didn't set the correct value.")

	t.Run("String", func(t *testing.T) {
		t.Run("positive", func(t *testing.T) {
			atom := NewInt32(math.MaxInt32)
			assert.Equal(t, "2147483647", atom.String(),
				"String() returned an unexpected value.")
		})

		t.Run("negative", func(t *testing.T) {
			atom := NewInt32(math.MinInt32)
			assert.Equal(t, "-2147483648", atom.String(),
				"String() returned an unexpected value.")
		})
	})
}
