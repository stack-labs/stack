package function

import (
	"github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/plugin"
	"github.com/stack-labs/stack-rpc/service"
)

var options struct {
	Stack struct {
		Service struct {
			Stack struct {
			} `sc:"stack"`
		} `sc:"service"`
	} `sc:"stack"`
}

type funcServicePlugin struct{}

func (s *funcServicePlugin) Name() string {
	return "stack"
}

func (s *funcServicePlugin) Options() []service.Option {
	return nil
}

func (s *funcServicePlugin) New(opts ...service.Option) service.Service {
	return NewFunction(opts...)
}

func init() {
	config.RegisterOptions(&options)
	plugin.ServicePlugins["function"] = &funcServicePlugin{}
}
