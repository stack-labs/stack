package stack

import (
	"github.com/stack-labs/stack/broker"
	"github.com/stack-labs/stack/broker/http"
)

type httpBrokerPlugin struct{}

func (h *httpBrokerPlugin) Name() string {
	return "http"
}

func (h *httpBrokerPlugin) Options() []broker.Option {
	return nil
}

func (h *httpBrokerPlugin) New(opts ...broker.Option) broker.Broker {
	return http.NewBroker(opts...)
}
