package grpc

import (
	"github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/plugin"
	"github.com/stack-labs/stack-rpc/service"
)

var options struct {
	Stack struct {
		Service struct {
			GRPC struct {
			} `sc:"grpc"`
		} `sc:"service"`
	} `sc:"stack"`
}

type grpcServicePlugin struct{}

func (g *grpcServicePlugin) Name() string {
	return "grpc"
}

func (g *grpcServicePlugin) Options() []service.Option {
	return nil
}

func (g *grpcServicePlugin) New(opts ...service.Option) service.Service {
	return NewService(opts...)
}

func init() {
	config.RegisterOptions(&options)
	plugin.ServicePlugins["grpc"] = &grpcServicePlugin{}
}
