// Package mucp initialises a mucp service
package mucp

import (
	// TODO: change to stack-rpc/service
	"github.com/stack-labs/stack-rpc"
	cmucp "github.com/stack-labs/stack-rpc/client/mucp"
	smucp "github.com/stack-labs/stack-rpc/server/mucp"
)

// NewService returns a new mucp service
func NewService(opts ...stack.Option) stack.Service {
	options := []stack.Option{
		stack.Client(cmucp.NewClient()),
		stack.Server(smucp.NewServer()),
	}

	options = append(options, opts...)

	return stack.NewService(options...)
}
