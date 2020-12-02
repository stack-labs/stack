package cmd

import (
	"context"

	"github.com/stack-labs/stack-rpc/broker"
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/registry"
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/transport"
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

	Broker    *broker.Broker
	Registry  *registry.Registry
	Selector  *selector.Selector
	Transport *transport.Transport
	Client    *client.Client
	Server    *server.Server
	Config    *config.Config
	Logger    *logger.Logger
}

// Command line Name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Command line Description
func Description(d string) Option {
	return func(o *Options) {
		o.Description = d
	}
}

// Command line Version
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

func Broker(b *broker.Broker) Option {
	return func(o *Options) {
		o.Broker = b
	}
}

func Selector(s *selector.Selector) Option {
	return func(o *Options) {
		o.Selector = s
	}
}

func Registry(r *registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

func Transport(t *transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}

func Client(c *client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

func Server(s *server.Server) Option {
	return func(o *Options) {
		o.Server = s
	}
}

func Config(c *config.Config) Option {
	return func(o *Options) {
		o.Config = c
	}
}

func Logger(log *logger.Logger) Option {
	return func(o *Options) {
		o.Logger = log
	}
}
