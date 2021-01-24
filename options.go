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
	"github.com/stack-labs/stack-rpc/service/web"
	"github.com/stack-labs/stack-rpc/transport"
)

func Cmd(c cmd.Cmd) service.Option {
	return service.Cmd(c)
}

func Logger(l logger.Logger) service.Option {
	return service.Logger(l)
}

func Broker(b broker.Broker) service.Option {
	return service.Broker(b)
}

func Client(c client.Client) service.Option {
	return service.Client(c)
}

func Config(c config.Config) service.Option {
	return service.Config(c)
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service.
// Can be used for extra service.Option values.
func Context(ctx context.Context) service.Option {
	return service.Context(ctx)
}

// HandleSignal toggles automatic installation of the signal handler that
// traps TERM, INT, and QUIT.  Users of this feature to disable the signal
// handler, should control liveness of the service through the context.
func HandleSignal(b bool) service.Option {
	return service.HandleSignal(b)
}

func Server(s server.Server) service.Option {
	return service.Server(s)
}

// Registry sets the registry for the service
// and the underlying components
func Registry(r registry.Registry) service.Option {
	return service.Registry(r)
}

// Selector sets the selector for the service client
func Selector(s selector.Selector) service.Option {
	return service.Selector(s)
}

// Transport sets the transport for the service
// and the underlying components
func Transport(t transport.Transport) service.Option {
	return service.Transport(t)
}

// Convenience service.Options

// Address sets the address of the server
func Address(addr string) service.Option {
	return service.Address(addr)
}

// Unique server id
func Id(id string) service.Option {
	return service.Id(id)
}

// Name of the service
func Name(n string) service.Option {
	return service.Name(n)
}

// Version of the service
func Version(v string) service.Option {
	return service.Version(v)
}

// Metadata associated with the service
func Metadata(md map[string]string) service.Option {
	return service.Metadata(md)
}

func WebRootPath(rootPath string) service.Option {
	return web.RootPath(rootPath)
}

func WebHandleFuncs(funcs ...web.HandlerFunc) service.Option {
	return web.HandleFuncs(funcs...)
}

func WebStaticDir(route, dir string) service.Option {
	return web.StaticDir(route, dir)
}

func Flags(flags ...cli.Flag) service.Option {
	return func(o *service.Options) {
		o.Cmd.App().Flags = append(o.Cmd.App().Flags, flags...)
	}
}

func Action(a func(*cli.Context)) service.Option {
	return func(o *service.Options) {
		o.Cmd.App().Action = a
	}
}

// Profile to be used for debug profile
func Profile(p profile.Profile) service.Option {
	return service.Profile(p)
}

// RegisterTTL specifies the TTL to use when registering the service
func RegisterTTL(t time.Duration) service.Option {
	return service.RegisterTTL(t)
}

// RegisterInterval specifies the interval on which to re-register
func RegisterInterval(t time.Duration) service.Option {
	return service.RegisterInterval(t)
}

// WrapClient is a convenience method for wrapping a Client with
// some middleware component. A list of wrappers can be provided.
// Wrappers are applied in reverse order so the last is executed first.
func WrapClient(w ...client.Wrapper) service.Option {
	return service.WrapClient(w...)
}

// WrapCall is a convenience method for wrapping a Client CallFunc
func WrapCall(w ...client.CallWrapper) service.Option {
	return service.WrapCall(w...)
}

// WrapHandler adds a handler Wrapper to a list of service.Options passed into the server
func WrapHandler(w ...server.HandlerWrapper) service.Option {
	return service.WrapHandler(w...)
}

// WrapSubscriber adds a subscriber Wrapper to a list of service.Options passed into the server
func WrapSubscriber(w ...server.SubscriberWrapper) service.Option {
	return service.WrapSubscriber(w...)
}

// Before and Afters

func BeforeStart(fn func() error) service.Option {
	return service.BeforeStart(fn)
}

func BeforeStop(fn func() error) service.Option {
	return service.BeforeStop(fn)
}

func AfterStart(fn func() error) service.Option {
	return service.AfterStart(fn)
}

func AfterStop(fn func() error) service.Option {
	return service.AfterStop(fn)
}
