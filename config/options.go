package config

import (
	"context"
)

type Options struct {
	EnableStorage bool
	StorageDir    string
	// for alternative data
	Context context.Context
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		EnableStorage: false,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

func EnableStorage(e bool) Option {
	return func(o *Options) {
		o.EnableStorage = e
	}
}

func StorageDir(d string) Option {
	return func(o *Options) {
		o.StorageDir = d
	}
}
