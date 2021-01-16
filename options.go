package stack

import (
	"context"
	"time"

	"github.com/stack-labs/stack-rpc/broker"
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/cmd"
	"github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/debug/profile"
	"github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/pkg/cli"
	"github.com/stack-labs/stack-rpc/registry"
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/service"
	"github.com/stack-labs/stack-rpc/transport"
)

type Options struct {
	Cmd         cmd.Cmd
	ServiceOpts []service.Option
}

func Cmd(c cmd.Cmd) Option {
	return func(o *Options) {
		o.Cmd = c
		o.ServiceOpts = append(o.ServiceOpts, service.Cmd(c))
	}
}

func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.Logger(l))
	}
}

func Broker(b broker.Broker) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.Broker(b))
	}
}

func Client(c client.Client) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.Client(c))
	}
}

func Config(c config.Config) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.Config(c))
	}
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service.
// Can be used for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.Context(ctx))
	}
}

// HandleSignal toggles automatic installation of the signal handler that
// traps TERM, INT, and QUIT.  Users of this feature to disable the signal
// handler, should control liveness of the service through the context.
func HandleSignal(b bool) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.HandleSignal(b))
	}
}

func Server(s server.Server) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.Server(s))
	}
}

// Registry sets the registry for the service
// and the underlying components
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, func(o *service.Options) {
			o.Registry = r
		})
	}
}

// Selector sets the selector for the service client
func Selector(s selector.Selector) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, func(o *service.Options) {
			o.Selector = s
		})
	}
}

// Transport sets the transport for the service
// and the underlying components
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, func(o *service.Options) {
			o.Transport = t
		})
	}
}

// Convenience options

// Address sets the address of the server
func Address(addr string) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.Address(addr))
	}
}

// Unique server id
func Id(id string) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.Id(id))
	}
}

// Name of the service
func Name(n string) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.Name(n))
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.Version(v))
	}
}

// Metadata associated with the service
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.Metadata(md))
	}
}

func Flags(flags ...cli.Flag) Option {
	return func(o *Options) {
		o.Cmd.App().Flags = append(o.Cmd.App().Flags, flags...)
	}
}

func Action(a func(*cli.Context)) Option {
	return func(o *Options) {
		o.Cmd.App().Action = a
	}
}

// Profile to be used for debug profile
func Profile(p profile.Profile) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.Profile(p))
	}
}

// RegisterTTL specifies the TTL to use when registering the service
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.RegisterTTL(t))
	}
}

// RegisterInterval specifies the interval on which to re-register
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.RegisterInterval(t))
	}
}

// WrapClient is a convenience method for wrapping a Client with
// some middleware component. A list of wrappers can be provided.
// Wrappers are applied in reverse order so the last is executed first.
func WrapClient(w ...client.Wrapper) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, func(o *service.Options) {
			o.ClientWrapper = w
		})
	}
}

// WrapCall is a convenience method for wrapping a Client CallFunc
func WrapCall(w ...client.CallWrapper) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.WrapCall(w...))
	}
}

// WrapHandler adds a handler Wrapper to a list of options passed into the server
func WrapHandler(w ...server.HandlerWrapper) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.WrapHandler(w...))
	}
}

// WrapSubscriber adds a subscriber Wrapper to a list of options passed into the server
func WrapSubscriber(w ...server.SubscriberWrapper) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.WrapSubscriber(w...))
	}
}

// Before and Afters

func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.BeforeStart(fn))
	}
}

func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.BeforeStop(fn))
	}
}

func AfterStart(fn func() error) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.AfterStart(fn))
	}
}

func AfterStop(fn func() error) Option {
	return func(o *Options) {
		o.ServiceOpts = append(o.ServiceOpts, service.AfterStop(fn))
	}
}
