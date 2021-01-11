package http

import (
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/plugin"
)

type httpClientPlugin struct {
}

func (h *httpClientPlugin) Name() string {
	return "http"
}

func (h *httpClientPlugin) Options() []client.Option {
	return nil
}

func (h *httpClientPlugin) New(opts ...client.Option) client.Client {
	return NewClient(opts...)
}

func init() {
	plugin.ClientPlugins["http"] = &httpClientPlugin{}
}
