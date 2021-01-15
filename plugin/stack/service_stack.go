package stack

import (
	"github.com/stack-labs/stack-rpc/service"
	"github.com/stack-labs/stack-rpc/service/stack"
)

type stackServicePlugin struct{}

func (s *stackServicePlugin) Name() string {
	return "stack"
}

func (s *stackServicePlugin) Options() []service.Option {
	return nil
}

func (s *stackServicePlugin) New(opts ...service.Option) service.Service {
	return stack.NewService(opts...)
}
