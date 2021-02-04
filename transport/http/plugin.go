package http

import (
	"github.com/stack-labs/stack/plugin"
	"github.com/stack-labs/stack/transport"
)

type httpTransportPlugin struct {
}

func (h *httpTransportPlugin) Name() string {
	return "http"
}

func (h *httpTransportPlugin) Options() []transport.Option {
	return nil
}

func (h *httpTransportPlugin) New(opts ...transport.Option) transport.Transport {
	return NewTransport(opts...)
}

func init() {
	plugin.TransportPlugins["http"] = &httpTransportPlugin{}
}
