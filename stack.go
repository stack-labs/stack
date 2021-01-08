// Package stack is a pluggable framework for microservices
package stack

import (
	"context"

	"github.com/stack-labs/stack-rpc/client"
	cmucp "github.com/stack-labs/stack-rpc/client/mucp"
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/service"
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
	return service.NewService(o.serviceOpts...)
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
	if c == nil {
		c = cmucp.NewClient()
	}
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
