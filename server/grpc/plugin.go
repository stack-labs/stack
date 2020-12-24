package grpc

import (
	"github.com/stack-labs/stack-rpc/plugin"
	"github.com/stack-labs/stack-rpc/server"
)

type grpcServerPlugin struct {
}

func (g *grpcServerPlugin) Name() string {
	return "grpc"
}

func (g *grpcServerPlugin) Options() []server.Option {
	return nil
}

func (g *grpcServerPlugin) New(opts ...server.Option) server.Server {
	return NewServer(opts...)
}

func init() {
	plugin.ServerPlugins["grpc"] = &grpcServerPlugin{}
}
