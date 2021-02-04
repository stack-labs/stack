package stack

import (
	"github.com/stack-labs/stack/registry"
	"github.com/stack-labs/stack/registry/mdns"
)

type mdnsRegistryPlugin struct{}

func (m *mdnsRegistryPlugin) Name() string {
	return "mdns"
}

func (m *mdnsRegistryPlugin) Options() []registry.Option {
	return nil
}

func (m *mdnsRegistryPlugin) New(opts ...registry.Option) registry.Registry {
	return mdns.NewRegistry(opts...)
}
