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

type stackServicePlugin struct{}

func (s *stackServicePlugin) Name() string {
	return "stack"
}

func (s *stackServicePlugin) Options() []service.Option {
	return nil
}

func (s *stackServicePlugin) New(opts ...service.Option) service.Service {
	return NewService(opts...)
}

func init() {
	config.RegisterOptions(&options)
	plugin.ServicePlugins["grpc"] = &stackServicePlugin{}
}
