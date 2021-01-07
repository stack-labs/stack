package service

import (
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/plugin"
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/util/log"
)

// Service is an interface that wraps the lower level libraries
// within stack-rpc. Its a convenience method for building
// and initialising services.
type Service interface {
	// The service name
	Name() string
	// Init initialises options
	Init(...Option) error
	// Options returns the current options
	Options() Options
	// Client is used to call services
	Client() client.Client
	// Server is for handling requests and events
	Server() server.Server
	//  Run the service
	Run() error
	// The service implementation
	String() string
}

func NewService(opts ...Option) Service {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}

	// todo check client is nil

	if len(o.ClientWrapper) > 0 {
		// apply in reverse
		for i := len(o.ClientWrapper); i > 0; i-- {
			o.Client = o.ClientWrapper[i-1](o.Client)
		}
	}

	if len(o.CallWrapper) > 0 {
		// todo move init client itself
		o.Client.Init(client.WrapCall(o.CallWrapper...))
	}

	// todo check server is nil

	if len(o.HandlerWrapper) > 0 {
		var wrappers []server.Option
		for _, wrap := range o.HandlerWrapper {
			wrappers = append(wrappers, server.WrapHandler(wrap))
		}
		// todo move init server itself
		o.Server.Init(wrappers...)
	}
	if len(o.SubscriberWrapper) > 0 {
		var wrappers []server.Option
		for _, wrap := range o.SubscriberWrapper {
			wrappers = append(wrappers, server.WrapSubscriber(wrap))
		}
		// todo move init server itself
		o.Server.Init(wrappers...)
	}

	p, ok := plugin.ServicePlugins[o.RPC]
	if !ok {
		log.Fatalf("[%s] service plugin isn't found")
	}

	return p.New(opts...)
}
