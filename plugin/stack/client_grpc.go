package stack

import (
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/client/grpc"
)

type grpcClientPlugin struct{}

func (m *grpcClientPlugin) Name() string {
	return "grpc"
}

func (m *grpcClientPlugin) Options() []client.Option {
	return nil
}

func (m *grpcClientPlugin) New(opts ...client.Option) client.Client {
	return grpc.NewClient(opts...)
}
