package memory

import (
	"github.com/stack-labs/stack-rpc/plugin"
	"github.com/stack-labs/stack-rpc/registry"
)

type memoryRegistryPlugin struct {
}

func (m *memoryRegistryPlugin) Name() string {
	return "memory"
}

func (m *memoryRegistryPlugin) Options() []registry.Option {
	return nil
}

func (m *memoryRegistryPlugin) New(opts ...registry.Option) registry.Registry {
	return NewRegistry(opts...)
}

func init() {
	plugin.RegistryPlugins["memory"] = &memoryRegistryPlugin{}
}
