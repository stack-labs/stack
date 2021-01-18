package cmd

import (
	"context"
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
