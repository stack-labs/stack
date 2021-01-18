package stack

import (
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/client/selector/static"
)

type staticSelectorPlugin struct{}

func (s *staticSelectorPlugin) Name() string {
	return "static"
}

func (s *staticSelectorPlugin) Options() []selector.Option {
	return nil
}

func (s *staticSelectorPlugin) New(opts ...selector.Option) selector.Selector {
	return static.NewSelector(opts...)
}
