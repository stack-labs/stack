package stack

import (
	"github.com/stack-labs/stack-rpc/registry"
	"github.com/stack-labs/stack-rpc/registry/mdns"
)

type mdnsRegistryPlugin struct{}

func (m *mdnsRegistryPlugin) Name() string {
	return "mdns"
}

func (m *mdnsRegistryPlugin) Options() []Registry.Option {
	return nil
}

func (m *mdnsRegistryPlugin) New(opts ...Registry.Option) Registry.Registry {
	return mdns.NewRegistry(opts...)
}
