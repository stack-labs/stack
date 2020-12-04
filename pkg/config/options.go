package config

import (
	"context"

	"github.com/stack-labs/stack-rpc/pkg/config/source"
)

type Options struct {
	Storage    bool
	StorageDir string
	Watch      bool
	Sources    []source.Source
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

func Sources(s ...source.Source) Option {
	return func(o *Options) {
		o.Sources = append(o.Sources, s...)
	}
}
