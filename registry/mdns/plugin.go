package mdns

import (
	"github.com/stack-labs/stack-rpc/plugin"
	"github.com/stack-labs/stack-rpc/registry"
)

type mdnsRegistryPlugin struct {
}

func (m *mdnsRegistryPlugin) Name() string {
	return "mdns"
}

func (m *mdnsRegistryPlugin) Options() []registry.Option {
	return nil
}

func (m *mdnsRegistryPlugin) New(opts ...registry.Option) registry.Registry {
	return NewRegistry(opts...)
}

func init() {
	plugin.RegistryPlugins["mdns"] = &mdnsRegistryPlugin{}
}
