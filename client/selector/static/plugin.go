package static

import (
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/plugin"
)

type staticSelectorPlugin struct {
}

func (s *staticSelectorPlugin) Name() string {
	return "static"
}

func (s *staticSelectorPlugin) Options() []selector.Option {
	return nil
}

func (s *staticSelectorPlugin) New(opts ...selector.Option) selector.Selector {
	return NewSelector(opts...)
}

func init() {
	plugin.SelectorPlugins["static"] = &staticSelectorPlugin{}
}
