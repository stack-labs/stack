// Package registry uses the stack-rpc registry for selection
package registry

import (
	"github.com/stack-labs/stack-rpc/client/selector"
)

// NewSelector returns a new registry selector
func NewSelector(opts ...selector.Option) selector.Selector {
	return selector.NewSelector(opts...)
}
