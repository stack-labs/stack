package stack

import (
	"github.com/stack-labs/stack-rpc/cmd"
	"github.com/stack-labs/stack-rpc/cmd/stack"
)

type cmdPlugin struct{}

func (h *cmdPlugin) Name() string {
	return "stack"
}

func (h *cmdPlugin) Options() []cmd.Option {
	return nil
}

func (h *cmdPlugin) New(opts ...cmd.Option) cmd.Cmd {
	return stack.NewCmd(opts...)
}

func init() {
	// plugin.CmdPlugins["stack"] = &cmdPlugin{}
}
