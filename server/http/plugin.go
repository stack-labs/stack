package http

import (
	"github.com/stack-labs/stack-rpc/plugin"
	"github.com/stack-labs/stack-rpc/server"
)

type httpServerPlugin struct{}

func (m *httpServerPlugin) Name() string {
	return "http"
}

func (m *httpServerPlugin) Options() []server.Option {
	return nil
}

func (m *httpServerPlugin) New(opts ...server.Option) server.Server {
	return NewServer(opts...)
}

func init() {
	plugin.ServerPlugins["http"] = &httpServerPlugin{}
}
