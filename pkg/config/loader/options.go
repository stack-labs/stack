package loader

import (
	"context"

	"github.com/stack-labs/stack-rpc/pkg/config/reader/json"

	"github.com/stack-labs/stack-rpc/pkg/config/reader"
)

type Options struct {
	Reader reader.Reader
	Watch  bool

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

// WithReader sets the config reader
func WithWatch(t bool) Option {
	return func(o *Options) {
		o.Watch = t
	}
}

func NewOptions(opts ...Option) Options {
	options := Options{
		Reader: json.NewReader(),
		Watch:  true,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}
