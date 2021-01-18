package stack

import (
	"github.com/stack-labs/stack-rpc/broker"
	"github.com/stack-labs/stack-rpc/broker/service"
)

type serviceBrokerPlugin struct{}

func (s *serviceBrokerPlugin) Name() string {
	return "service"
}

func (s *serviceBrokerPlugin) Options() []broker.Option {
	return nil
}

func (s *serviceBrokerPlugin) New(opts ...broker.Option) broker.Broker {
	return service.NewBroker(opts...)
}
