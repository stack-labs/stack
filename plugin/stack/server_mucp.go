package stack

import (
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/server/mucp"
)

type mucpServerPlugin struct {}

func (m *mucpServerPlugin) Name() string {
	return "mucp"
}

func (m *mucpServerPlugin) Options() []server.Option {
	return nil
}

func (m *mucpServerPlugin) New(opts ...server.Option) server.Server {
	return mucp.NewServer(opts...)
}

func init() {

}
