package stack

import (
	"github.com/stack-labs/stack-rpc/service"
	"github.com/stack-labs/stack-rpc/service/grpc"
)

type grpcServicePlugin struct{}

func (g *grpcServicePlugin) Name() string {
	return "grpc"
}

func (g *grpcServicePlugin) Options() []service.Option {
	return nil
}

func (g *grpcServicePlugin) New(opts ...service.Option) service.Service {
	return grpc.NewService(opts...)
}
