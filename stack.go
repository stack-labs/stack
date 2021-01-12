// Package stack is a pluggable framework for microservices
package stack

import (
	"context"

	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/plugin"
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/service"
	"github.com/stack-labs/stack-rpc/util/log"

	_ "github.com/stack-labs/stack-rpc/broker/http"
	_ "github.com/stack-labs/stack-rpc/client/mucp"
	_ "github.com/stack-labs/stack-rpc/logger/console"
	_ "github.com/stack-labs/stack-rpc/registry/mdns"
	_ "github.com/stack-labs/stack-rpc/server/mucp"
	_ "github.com/stack-labs/stack-rpc/service/grpc"
	_ "github.com/stack-labs/stack-rpc/service/stack"
	_ "github.com/stack-labs/stack-rpc/service/web"
	_ "github.com/stack-labs/stack-rpc/transport/http"
)

type serviceKey struct{}

// Publisher is syntactic sugar for publishing
type Publisher interface {
	Publish(ctx context.Context, msg interface{}, opts ...client.PublishOption) error
}

type Option func(*Options)

// NewService creates and returns a new Service based on the packages within.
func NewService(opts ...Option) service.Service {
	o := newOptions(opts...)

	// set default
	// this will be removed in future
	so := &service.Options{}
	for _, opt := range o.ServiceOpts {
		opt(so)
	}

	p, ok := plugin.ServicePlugins[so.RPC]
	if !ok {
		log.Fatal("[%s] service plugin isn't found", so.RPC)
	}

	return p.New(o.ServiceOpts...)
}

// FromContext retrieves a Service from the Context.
func FromContext(ctx context.Context) (service.Service, bool) {
	s, ok := ctx.Value(serviceKey{}).(service.Service)
	return s, ok
}

// NewContext returns a new Context with the Service embedded within it.
func NewContext(ctx context.Context, s service.Service) context.Context {
	return context.WithValue(ctx, serviceKey{}, s)
}

// NewPublisher returns a new Publisher
func NewPublisher(topic string, c client.Client) Publisher {
	return &publisher{c, topic}
}

// RegisterHandler is syntactic sugar for registering a handler
func RegisterHandler(s server.Server, h interface{}, opts ...server.HandlerOption) error {
	return s.Handle(s.NewHandler(h, opts...))
}

// RegisterSubscriber is syntactic sugar for registering a subscriber
func RegisterSubscriber(topic string, s server.Server, h interface{}, opts ...server.SubscriberOption) error {
	return s.Subscribe(s.NewSubscriber(topic, h, opts...))
}
