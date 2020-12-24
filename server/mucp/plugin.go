package mucp

import (
	"github.com/stack-labs/stack-rpc/plugin"
	"github.com/stack-labs/stack-rpc/server"
)

type mucpServerPlugin struct {
}

func (m *mucpServerPlugin) Name() string {
	return "mucp"
}

func (m *mucpServerPlugin) Options() []server.Option {
	return nil
}

func (m *mucpServerPlugin) New(opts ...server.Option) server.Server {
	return NewServer(opts...)
}

func init() {
	plugin.ServerPlugins["mucp"] = &mucpServerPlugin{}
}
