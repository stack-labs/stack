package dns

import (
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/plugin"
)

type dnsSelectorPlugin struct {
}

func (h *dnsSelectorPlugin) Name() string {
	return "dns"
}

func (h *dnsSelectorPlugin) Options() []selector.Option {
	return nil
}

func (h *dnsSelectorPlugin) New(opts ...selector.Option) selector.Selector {
	return NewSelector(opts...)
}

func init() {
	plugin.SelectorPlugins["dns"] = &dnsSelectorPlugin{}
}
