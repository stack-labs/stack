// Package mucp provides an mucp client
package mucp

import (
	"github.com/stack-labs/stack-rpc/client"
)

// NewClient returns a new stack client interface
func NewClient(opts ...client.Option) client.Client {
	return client.NewClient(opts...)
}
