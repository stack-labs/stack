package stack

import (
	"github.com/stack-labs/stack-rpc/registry"
	"github.com/stack-labs/stack-rpc/registry/service"
)

type serviceRegistryPlugin struct{}

func (s *serviceRegistryPlugin) Name() string {
	return "service"
}

func (s *serviceRegistryPlugin) Options() []Registry.Option {
	return nil
}

func (s *serviceRegistryPlugin) New(opts ...Registry.Option) Registry.Registry {
	return service.NewRegistry(opts...)
}
