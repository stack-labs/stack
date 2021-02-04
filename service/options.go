package service

import (
	"context"
	"time"

	"github.com/stack-labs/stack/auth"
	"github.com/stack-labs/stack/broker"
	"github.com/stack-labs/stack/client"
	"github.com/stack-labs/stack/client/selector"
	"github.com/stack-labs/stack/cmd"
	"github.com/stack-labs/stack/config"
	"github.com/stack-labs/stack/debug/profile"
	"github.com/stack-labs/stack/logger"
	"github.com/stack-labs/stack/registry"
	"github.com/stack-labs/stack/server"
	"github.com/stack-labs/stack/transport"
)

type Option func(o *Options)

type Options struct {
	// maybe put them in metadata is better
	Id   string
	Name string
	RPC  string
	Cmd  cmd.Cmd
	Conf string

	BrokerOptions    BrokerOptions
	ClientOptions    ClientOptions
	ServerOptions    ServerOptions
	RegistryOptions  RegistryOptions
	TransportOptions TransportOptions
	SelectorOptions  SelectorOptions
	ConfigOptions    ConfigOptions
	LoggerOptions    LoggerOptions
	AuthOptions      AuthOptions

	Broker    broker.Broker
	Client    client.Client
	Server    server.Server
	Registry  registry.Registry
	Transport transport.Transport
	Selector  selector.Selector
	Config    config.Config
	Logger    logger.Logger
	Auth      auth.Auth
	Profile   profile.Profile

	// Before and After funcs
	BeforeInit  []func(sOpts *Options) error
	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error

	ClientWrapper     []client.Wrapper
	CallWrapper       []client.CallWrapper
	HandlerWrapper    []server.HandlerWrapper
	SubscriberWrapper []server.SubscriberWrapper
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context

	Signal bool
}

type BrokerOptions []broker.Option

func (b BrokerOptions) Options() broker.Options {
	opts := broker.Options{}
	for _, o := range b {
		o(&opts)
	}

	return opts
}

type ClientOptions []client.Option

func (c ClientOptions) Options() client.Options {
	opts := client.Options{}
	for _, o := range c {
		o(&opts)
	}

	return opts
}

type ServerOptions []server.Option

func (c ServerOptions) Options() server.Options {
	opts := server.Options{}
	for _, o := range c {
		o(&opts)
	}

	return opts
}

type RegistryOptions []registry.Option

func (c RegistryOptions) Options() registry.Options {
	opts := registry.Options{}
	for _, o := range c {
		o(&opts)
	}

	return opts
}

type TransportOptions []transport.Option

func (c TransportOptions) Options() transport.Options {
	opts := transport.Options{}
	for _, o := range c {
		o(&opts)
	}

	return opts
}

type SelectorOptions []selector.Option

func (c SelectorOptions) Options() selector.Options {
	opts := selector.Options{}
	for _, o := range c {
		o(&opts)
	}

	return opts
}

type ConfigOptions []config.Option

func (c ConfigOptions) Options() config.Options {
	opts := config.Options{}
	for _, o := range c {
		o(&opts)
	}

	return opts
}

type LoggerOptions []logger.Option

func (c LoggerOptions) Options() logger.Options {
	opts := logger.Options{}
	for _, o := range c {
		o(&opts)
	}

	return opts
}

type AuthOptions []auth.Option

func (a AuthOptions) Options() auth.Options {
	opts := auth.Options{}
	for _, o := range a {
		o(&opts)
	}

	return opts
}

func Cmd(c cmd.Cmd) Option {
	return func(o *Options) {
		o.Cmd = c
	}
}

// RPC sets the type of service, eg. stack, grpc
// but this func will be deprecated
func RPC(r string) Option {
	return func(o *Options) {
		o.RPC = r
	}
}

func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

func Broker(b broker.Broker) Option {
	return func(o *Options) {
		o.Broker = b
	}
}

func Client(c client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

func Config(c config.Config) Option {
	return func(o *Options) {
		o.Config = c
	}
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service.
// Can be used for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// HandleSignal toggles automatic installation of the signal handler that
// traps TERM, INT, and QUIT.  Users of this feature to disable the signal
// handler, should control liveness of the service through the context.
func HandleSignal(b bool) Option {
	return func(o *Options) {
		o.Signal = b
	}
}

func Server(s server.Server) Option {
	return func(o *Options) {
		o.Server = s
	}
}

// Registry sets the registry for the service
// and the underlying components
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// Selector sets the selector for the service client
func Selector(s selector.Selector) Option {
	return func(o *Options) {
		o.Selector = s
	}
}

// Transport sets the transport for the service
// and the underlying components
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}

// Address sets the address of the server
func Address(addr string) Option {
	return func(o *Options) {
		o.ServerOptions = append(o.ServerOptions, server.Address(addr))
	}
}

// Unique server id
func Id(id string) Option {
	return func(o *Options) {
		o.Id = id
		o.ServerOptions = append(o.ServerOptions, server.Id(id))
	}
}

// Name of the service
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
		o.ServerOptions = append(o.ServerOptions, server.Name(n))
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		o.ServerOptions = append(o.ServerOptions, server.Version(v))
	}
}

// Metadata associated with the service
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.ServerOptions = append(o.ServerOptions, server.Metadata(md))
	}
}

// Profile to be used for debug profile
func Profile(p profile.Profile) Option {
	return func(o *Options) {
		o.Profile = p
	}
}

func Auth(au auth.Auth) Option {
	return func(o *Options) {
		o.Auth = au
	}
}

// RegisterTTL specifies the TTL to use when registering the service
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		o.ServerOptions = append(o.ServerOptions, server.RegisterTTL(t))
	}
}

// RegisterInterval specifies the interval on which to re-register
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		o.ServerOptions = append(o.ServerOptions, server.RegisterInterval(t))
	}
}

// WrapClient is a convenience method for wrapping a Client with
// some middleware component. A list of wrappers can be provided.
// Wrappers are applied in reverse order so the last is executed first.
func WrapClient(w ...client.Wrapper) Option {
	return func(o *Options) {
		o.ClientWrapper = w
	}
}

// WrapCall is a convenience method for wrapping a Client CallFunc
func WrapCall(w ...client.CallWrapper) Option {
	return func(o *Options) {
		o.CallWrapper = w
	}
}

// WrapHandler adds a handler Wrapper to a list of options passed into the server
func WrapHandler(w ...server.HandlerWrapper) Option {
	return func(o *Options) {
		o.HandlerWrapper = w
	}
}

// WrapSubscriber adds a subscriber Wrapper to a list of options passed into the server
func WrapSubscriber(w ...server.SubscriberWrapper) Option {
	return func(o *Options) {
		o.SubscriberWrapper = w
	}
}

func BeforeInit(fn func(sOpts *Options) error) Option {
	return func(o *Options) {
		o.BeforeInit = append(o.BeforeInit, fn)
	}
}

func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

func AfterStart(fn func() error) Option {
	return func(o *Options) {
		o.AfterStart = append(o.AfterStart, fn)
	}
}

func AfterStop(fn func() error) Option {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}
