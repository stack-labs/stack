package stack

import (
	"github.com/stack-labs/stack-rpc/registry"
	"github.com/stack-labs/stack-rpc/registry/memory"
)

type memoryRegistryPlugin struct{}

func (m *memoryRegistryPlugin) Name() string {
	return "memory"
}

func (m *memoryRegistryPlugin) Options() []Registry.Option {
	return nil
}

func (m *memoryRegistryPlugin) New(opts ...Registry.Option) Registry.Registry {
	return memory.NewRegistry(opts...)
}
