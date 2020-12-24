package registry

import (
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/plugin"
)

type cacheSelectorPlugin struct {
}

func (c *cacheSelectorPlugin) Name() string {
	return "cache"
}

func (c *cacheSelectorPlugin) Options() []selector.Option {
	return nil
}

func (c *cacheSelectorPlugin) New(opts ...selector.Option) selector.Selector {
	return NewSelector(opts...)
}

func init() {
	plugin.SelectorPlugins["cache"] = &cacheSelectorPlugin{}
}
