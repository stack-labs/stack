package mucp

import (
	"context"

	"github.com/stack-labs/stack-rpc/server"
)

type serverKey struct{}

func FromContext(ctx context.Context) (server.Server, bool) {
	c, ok := ctx.Value(serverKey{}).(server.Server)
	return c, ok
}

func NewContext(ctx context.Context, s server.Server) context.Context {
	return context.WithValue(ctx, serverKey{}, s)
}
