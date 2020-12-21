// Package selector is a way to pick a list of service nodes
package selector

import (
	"errors"

	"github.com/stack-labs/stack-rpc/registry"
)

// Selector builds on the registry as a mechanism to pick nodes
// and mark their status. This allows host pools and other things
// to be built using various algorithms.
type Selector interface {
	Init(opts ...Option) error
	Options() Options
	// Next returns the next node
	Next(service string, opts ...SelectOption) (*registry.Node, error)
	// Mark sets the success/error against a node
	Mark(service string, node *registry.Node, err error)
	// Reset returns state back to zero for a service
	Reset(service string)
	// Close renders the selector unusable
	Close() error
	// Name of the selector
	String() string
}

var (
	ErrNotFound      = errors.New("not found")
	ErrNoneAvailable = errors.New("none available")
)
