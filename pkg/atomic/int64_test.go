package atomic

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt64(t *testing.T) {
	atom := NewInt64(42)

	require.Equal(t, int64(42), atom.Load(), "Load didn't work.")
	require.Equal(t, int64(46), atom.Add(4), "Add didn't work.")
	require.Equal(t, int64(44), atom.Sub(2), "Sub didn't work.")
	require.Equal(t, int64(45), atom.Inc(), "Inc didn't work.")
	require.Equal(t, int64(44), atom.Dec(), "Dec didn't work.")

	require.True(t, atom.CAS(44, 0), "CAS didn't report a swap.")
	require.Equal(t, int64(0), atom.Load(), "CAS didn't set the correct value.")

	require.Equal(t, int64(0), atom.Swap(1), "Swap didn't return the old value.")
	require.Equal(t, int64(1), atom.Load(), "Swap didn't set the correct value.")

	atom.Store(42)
	require.Equal(t, int64(42), atom.Load(), "Store didn't set the correct value.")

	t.Run("String", func(t *testing.T) {
		t.Run("positive", func(t *testing.T) {
			atom := NewInt64(math.MaxInt64)
			assert.Equal(t, "9223372036854775807", atom.String(),
				"String() returned an unexpected value.")
		})

		t.Run("negative", func(t *testing.T) {
			atom := NewInt64(math.MinInt64)
			assert.Equal(t, "-9223372036854775808", atom.String(),
				"String() returned an unexpected value.")
		})
	})
}
