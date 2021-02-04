package stack

import (
	"github.com/stack-labs/stack/server"
	"github.com/stack-labs/stack/server/mock"
)

type mockServerPlugin struct{}

func (m *mockServerPlugin) Name() string {
	return "mock"
}

func (m *mockServerPlugin) Options() []server.Option {
	return nil
}

func (m *mockServerPlugin) New(opts ...server.Option) server.Server {
	return mock.NewServer(opts...)
}
