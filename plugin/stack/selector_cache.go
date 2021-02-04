package stack

import (
	"github.com/stack-labs/stack/client/selector"
	"github.com/stack-labs/stack/client/selector/registry"
)

type cacheSelectorPlugin struct{}

func (c *cacheSelectorPlugin) Name() string {
	return "cache"
}

func (c *cacheSelectorPlugin) Options() []selector.Option {
	return nil
}

func (c *cacheSelectorPlugin) New(opts ...selector.Option) selector.Selector {
	return registry.NewSelector(opts...)
}
