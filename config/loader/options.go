package loader

import (
	"context"

	"github.com/stack-labs/stack-rpc/config/reader/json"

	"github.com/stack-labs/stack-rpc/config/reader"
)

type Options struct {
	Reader reader.Reader

	// for alternative data
	Context context.Context
}

type Option func(o *Options)

// WithReader sets the config reader
func WithReader(r reader.Reader) Option {
	return func(o *Options) {
		o.Reader = r
	}
}

func NewOptions(opts ...Option) Options {
	options := Options{
		Reader: json.NewReader(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}
