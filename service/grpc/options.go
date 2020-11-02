package grpc

import (
	"crypto/tls"

	"github.com/stack-labs/stack-rpc"
	gc "github.com/stack-labs/stack-rpc/client/grpc"
	gs "github.com/stack-labs/stack-rpc/server/grpc"
)

// WithTLS sets the TLS config for the service
func WithTLS(t *tls.Config) stack.Option {
	return func(o *stack.Options) {
		o.Client.Init(
			gc.AuthTLS(t),
		)
		o.Server.Init(
			gs.AuthTLS(t),
		)
	}
}
