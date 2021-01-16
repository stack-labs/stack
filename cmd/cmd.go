// Package cmd is an interface for parsing the command line
package cmd

import (
	"github.com/stack-labs/stack-rpc/pkg/cli"
	"github.com/stack-labs/stack-rpc/plugin"
)

type Cmd interface {
	// The cli app within this cmd
	App() *cli.App
	// Adds options, parses flags and initialise
	// exits on error
	Init(opts ...Option) error
	// Options set within this command
	Options() Options
	// ConfigFile path. This is not good
	ConfigFile() string
}

func NewCmd(opts ...Option) Cmd {
	// todo optional
	return plugin.CmdPlugins["stack"].New(opts...)
}
