package cmd

import (
	"context"

	au "github.com/stack-labs/stack-rpc/auth"
	bk "github.com/stack-labs/stack-rpc/broker"
	cl "github.com/stack-labs/stack-rpc/client"
	sel "github.com/stack-labs/stack-rpc/client/selector"
	cfg "github.com/stack-labs/stack-rpc/config"
	log "github.com/stack-labs/stack-rpc/logger"
	reg "github.com/stack-labs/stack-rpc/registry"
	ser "github.com/stack-labs/stack-rpc/server"
	tra "github.com/stack-labs/stack-rpc/transport"
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

	Broker    *bk.Broker
	Registry  *reg.Registry
	Selector  *sel.Selector
	Transport *tra.Transport
	Client    *cl.Client
	Server    *ser.Server
	Config    *cfg.Config
	Logger    *log.Logger
	Auth      *au.Auth
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

func Broker(b *bk.Broker) Option {
	return func(o *Options) {
		o.Broker = b
	}
}

func Selector(s *sel.Selector) Option {
	return func(o *Options) {
		o.Selector = s
	}
}

func Registry(r *reg.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

func Transport(t *tra.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}

func Client(c *cl.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

func Server(s *ser.Server) Option {
	return func(o *Options) {
		o.Server = s
	}
}

func Config(c *cfg.Config) Option {
	return func(o *Options) {
		o.Config = c
	}
}

func Auth(a *au.Auth) Option {
	return func(o *Options) {
		o.Auth = a
	}
}

func Logger(log *log.Logger) Option {
	return func(o *Options) {
		o.Logger = log
	}
}
