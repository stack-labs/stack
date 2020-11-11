package handler

import (
	"github.com/stack-labs/stack-rpc"
	"github.com/stack-labs/stack-rpc/api/router"
)

type Options struct {
	Namespace string
	Router    router.Router
	Service   stack.Service
}

type Option func(o *Options)

// NewOptions fills in the blanks
func NewOptions(opts ...Option) Options {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	// create service if its blank
	if options.Service == nil {
		WithService(stack.NewService())(&options)
	}

	// set namespace if blank
	if len(options.Namespace) == 0 {
		WithNamespace("stack.rpc.api")(&options)
	}

	return options
}

// WithNamespace specifies the namespace for the handler
func WithNamespace(s string) Option {
	return func(o *Options) {
		o.Namespace = s
	}
}

// WithRouter specifies a router to be used by the handler
func WithRouter(r router.Router) Option {
	return func(o *Options) {
		o.Router = r
	}
}

// WithService specifies a stack.Service
func WithService(s stack.Service) Option {
	return func(o *Options) {
		o.Service = s
	}
}
