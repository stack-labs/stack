package config

import (
	"context"
)

type Options struct {
	Storage    bool
	StorageDir string
	Watch      bool
	// for alternative data
	Context context.Context
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Storage: false,
		Watch:   true,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

func Watch(t bool) Option {
	return func(o *Options) {
		o.Watch = t
	}
}

func Storage(e bool) Option {
	return func(o *Options) {
		o.Storage = e
	}
}

func StorageDir(d string) Option {
	return func(o *Options) {
		o.StorageDir = d
	}
}
