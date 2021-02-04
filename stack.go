// Package stack is a pluggable framework for microservices
package stack

import (
	"context"

	"github.com/stack-labs/stack/client"
	"github.com/stack-labs/stack/server"
	"github.com/stack-labs/stack/service"
	"github.com/stack-labs/stack/service/grpc"
	"github.com/stack-labs/stack/service/stack"
	"github.com/stack-labs/stack/service/web"

	_ "github.com/stack-labs/stack/plugin/stack"
)

type serviceKey struct{}

// Publisher is syntactic sugar for publishing
type Publisher interface {
	Publish(ctx context.Context, msg interface{}, opts ...client.PublishOption) error
}

// NewService creates and returns a new Service based on the packages within.
func NewService(opts ...service.Option) service.Service {
	return stack.NewService(opts...)
}

// NewWebService creates and returns a new web Service based on the packages within.
func NewWebService(opts ...service.Option) service.Service {
	return stack.NewService(web.NewOptions(opts...)...)
}

// NewGRPCService creates and returns a new web Service based on the packages within.
func NewGRPCService(opts ...service.Option) service.Service {
	return stack.NewService(grpc.NewOptions(opts...)...)
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
