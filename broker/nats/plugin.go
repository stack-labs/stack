package nats

import (
	"github.com/stack-labs/stack-rpc/broker"
	"github.com/stack-labs/stack-rpc/plugin"
)

type natsBrokerPlugin struct {
}

func (n *natsBrokerPlugin) Name() string {
	return "nats"
}

func (n *natsBrokerPlugin) Options() []broker.Option {
	return nil
}

func (n *natsBrokerPlugin) New(opts ...broker.Option) broker.Broker {
	return NewBroker(opts...)
}

func init() {
	plugin.BrokerPlugins["nats"] = &natsBrokerPlugin{}
}
