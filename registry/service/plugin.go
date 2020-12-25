package service

import (
	"github.com/stack-labs/stack-rpc/plugin"
	"github.com/stack-labs/stack-rpc/registry"
)

type serviceRegistryPlugin struct {
}

func (s *serviceRegistryPlugin) Name() string {
	return "http"
}

func (s *serviceRegistryPlugin) Options() []registry.Option {
	return nil
}

func (s *serviceRegistryPlugin) New(opts ...registry.Option) registry.Registry {
	return NewRegistry(opts...)
}

func init() {
	plugin.RegistryPlugins["service"] = &serviceRegistryPlugin{}
}
