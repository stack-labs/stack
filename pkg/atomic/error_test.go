package atomic

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorByValue(t *testing.T) {
	err := &Error{}
	require.Nil(t, err.Load(), "Initial value shall be nil")
}

func TestNewErrorWithNilArgument(t *testing.T) {
	err := NewError(nil)
	require.Nil(t, err.Load(), "Initial value shall be nil")
}

func TestErrorCanStoreNil(t *testing.T) {
	err := NewError(errors.New("hello"))
	err.Store(nil)
	require.Nil(t, err.Load(), "Stored value shall be nil")
}

func TestNewErrorWithError(t *testing.T) {
	err1 := errors.New("hello1")
	err2 := errors.New("hello2")

	atom := NewError(err1)
	require.Equal(t, err1, atom.Load(), "Expected Load to return initialized value")

	atom.Store(err2)
	require.Equal(t, err2, atom.Load(), "Expected Load to return overridden value")
}
