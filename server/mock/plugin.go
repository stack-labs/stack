package mock

import (
	"github.com/stack-labs/stack-rpc/plugin"
	"github.com/stack-labs/stack-rpc/server"
)

type mockServerPlugin struct {
}

func (m *mockServerPlugin) Name() string {
	return "mock"
}

func (m *mockServerPlugin) Options() []server.Option {
	return nil
}

func (m *mockServerPlugin) New(opts ...server.Option) server.Server {
	return NewServer(opts...)
}

func init() {
	plugin.ServerPlugins["mock"] = &mockServerPlugin{}
}
