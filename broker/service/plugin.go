package service

import (
	"github.com/stack-labs/stack-rpc/broker"
	"github.com/stack-labs/stack-rpc/plugin"
)

type serviceBrokerPlugin struct {
}

func (s *serviceBrokerPlugin) Name() string {
	return "service"
}

func (s *serviceBrokerPlugin) Options() []broker.Option {
	return nil
}

func (s *serviceBrokerPlugin) New(opts ...broker.Option) broker.Broker {
	return NewBroker(opts...)
}

func init() {
	plugin.BrokerPlugins["service"] = &serviceBrokerPlugin{}
}
