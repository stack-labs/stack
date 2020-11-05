package config

import (
	"context"

	"github.com/stack-labs/stack-rpc/config/loader"
	"github.com/stack-labs/stack-rpc/config/loader/memory"
	"github.com/stack-labs/stack-rpc/config/reader"
	"github.com/stack-labs/stack-rpc/config/reader/json"
)

type Options struct {
	Loader loader.Loader
	Reader reader.Reader

	EnableStorage bool
	// for alternative data
	Context context.Context
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Loader:        memory.NewLoader(),
		Reader:        json.NewReader(),
		EnableStorage: false,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// Loader sets the loader for manager config
func Loader(l loader.Loader) Option {
	return func(o *Options) {
		o.Loader = l
	}
}

// Reader sets the config reader
func Reader(r reader.Reader) Option {
	return func(o *Options) {
		o.Reader = r
	}
}

func EnableStorage(e bool) Option {
	return func(o *Options) {
		o.EnableStorage = e
	}
}
