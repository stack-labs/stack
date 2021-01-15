package stack

import (
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/server/http"
)

type httpServerPlugin struct{}

func (m *httpServerPlugin) Name() string {
	return "http"
}

func (m *httpServerPlugin) Options() []server.Option {
	return nil
}

func (m *httpServerPlugin) New(opts ...server.Option) server.Server {
	return http.NewServer(opts...)
}
