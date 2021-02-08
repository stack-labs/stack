package stack

import (
	"github.com/stack-labs/stack/server"
	"github.com/stack-labs/stack/server/grpc"
)

type grpcServerPlugin struct{}

func (g *grpcServerPlugin) Name() string {
	return "grpc"
}

func (g *grpcServerPlugin) Options() []server.Option {
	return nil
}

func (g *grpcServerPlugin) New(opts ...server.Option) server.Server {
	return grpc.NewServer(opts...)
}
