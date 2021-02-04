package zookeeper

import (
	"github.com/stack-labs/stack/plugin"
	"github.com/stack-labs/stack/registry"
)

type zkRegistryPlugin struct {
}

func (z *zkRegistryPlugin) Name() string {
	return "zookeeper"
}

func (z *zkRegistryPlugin) Options() []registry.Option {
	return nil
}

func (z *zkRegistryPlugin) New(opts ...registry.Option) registry.Registry {
	return NewRegistry(opts...)
}

func init() {
	plugin.RegistryPlugins["zookeeper"] = &zkRegistryPlugin{}
}
