package memory

import (
	"github.com/stack-labs/stack-rpc/broker"
	"github.com/stack-labs/stack-rpc/plugin"
)

type memoryBrokerPlugin struct {
}

func (m *memoryBrokerPlugin) Name() string {
	return "memory"
}

func (m *memoryBrokerPlugin) Options() []broker.Option {
	return nil
}

func (m *memoryBrokerPlugin) New(opts ...broker.Option) broker.Broker {
	return NewBroker(opts...)
}

func init() {
	plugin.BrokerPlugins["memory"] = &memoryBrokerPlugin{}
}
