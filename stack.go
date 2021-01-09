// Package stack is a pluggable framework for microservices
package stack

import (
	"context"
	"github.com/stack-labs/stack-rpc/util/log"

	"github.com/stack-labs/stack-rpc/client"
	cmucp "github.com/stack-labs/stack-rpc/client/mucp"
	"github.com/stack-labs/stack-rpc/cmd"
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

	o.serviceOpts = append(o.serviceOpts, service.BeforeInit(func(sOpts *service.Options) error {
		// cmd helps stack parse command options and reset the options that should work.
		if err := o.Cmd.Init(
			cmd.Broker(&sOpts.Broker),
			cmd.Registry(&sOpts.Registry),
			cmd.Transport(&sOpts.Transport),
			cmd.Client(&sOpts.Client),
			cmd.Server(&sOpts.Server),
			cmd.Selector(&sOpts.Selector),
			cmd.Logger(&sOpts.Logger),
			cmd.Config(&sOpts.Config),
			cmd.Auth(&sOpts.Auth),
		); err != nil {
			log.Errorf("cmd init error: %s", err)
			return err
		}

		return nil
	}))

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
