package service

import (
	"context"
	"github.com/stack-labs/stack-rpc/broker/http"
	clientM "github.com/stack-labs/stack-rpc/client/mucp"
	selectorR "github.com/stack-labs/stack-rpc/client/selector/registry"
	"github.com/stack-labs/stack-rpc/registry/mdns"
	"github.com/stack-labs/stack-rpc/server/mucp"
	transportH "github.com/stack-labs/stack-rpc/transport/http"

	"github.com/stack-labs/stack-rpc/auth"
	"github.com/stack-labs/stack-rpc/broker"
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/debug/profile"
	"github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/registry"
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/transport"
)

type Option func(o *Options)

type Options struct {
	// maybe put them in metadata is better
	Id   string
	Name string
	RPC  string

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
	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context

	Signal bool
}

func NewOptions(opts ...Option) Options {
	opt := Options{
		Broker:    http.NewBroker(),
		Client:    clientM.NewClient(),
		Server:    mucp.NewServer(),
		Registry:  mdns.NewRegistry(),
		Transport: transportH.NewTransport(),
		Selector:  selectorR.NewSelector(),
		Logger:    logger.DefaultLogger,
		Config:    config.DefaultConfig,
		Auth:      auth.NoopAuth,
		Context:   context.Background(),
		Signal:    true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}
