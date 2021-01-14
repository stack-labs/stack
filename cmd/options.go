package cmd

import (
	"context"

	cfg "github.com/stack-labs/stack-rpc/config"
	ss "github.com/stack-labs/stack-rpc/service"
)

type Option func(o *Options)

type Options struct {
	// For the Command Line itself
	Name        string
	Description string
	Version     string
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context

	ServiceOptions *ss.Options
	Config         *cfg.Config
}

func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

func Description(d string) Option {
	return func(o *Options) {
		o.Description = d
	}
}

func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

func ServiceOptions(so *ss.Options) Option {
	return func(o *Options) {
		o.ServiceOptions = so
	}
}

func Config(c *cfg.Config) Option {
	return func(o *Options) {
		o.Config = c
	}
}