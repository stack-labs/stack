package stack

import (
	"github.com/stack-labs/stack/client"
	"github.com/stack-labs/stack/client/http"
	"github.com/stack-labs/stack/plugin"
)

type httpClientPlugin struct {}

func (h *httpClientPlugin) Name() string {
	return "http"
}

func (h *httpClientPlugin) Options() []client.Option {
	return nil
}

func (h *httpClientPlugin) New(opts ...client.Option) client.Client {
	return http.NewClient(opts...)
}

func init() {
	plugin.ClientPlugins["http"] = &httpClientPlugin{}
}
