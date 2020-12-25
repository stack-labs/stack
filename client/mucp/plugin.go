package mucp

import (
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/plugin"
)

type mucpClientPlugin struct {
}

func (m *mucpClientPlugin) Name() string {
	return "mucp"
}

func (m *mucpClientPlugin) Options() []client.Option {
	return nil
}

func (m *mucpClientPlugin) New(opts ...client.Option) client.Client {
	return NewClient(opts...)
}

func init() {
	plugin.ClientPlugins["mucp"] = &mucpClientPlugin{}
}
