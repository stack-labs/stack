package grpc

import (
	"crypto/tls"

	gc "github.com/stack-labs/stack-rpc/client/grpc"
	gs "github.com/stack-labs/stack-rpc/server/grpc"
	"github.com/stack-labs/stack-rpc/service"
)

// WithTLS sets the TLS config for the service
func WithTLS(t *tls.Config) service.Option {
	return func(o *service.Options) {
		o.Client.Init(
			gc.AuthTLS(t),
		)
		o.Server.Init(
			gs.AuthTLS(t),
		)
	}
}
