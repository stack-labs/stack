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

	p, ok := plugin.ServicePlugins[o.RPC]
	if !ok {
		log.Fatalf("[%s] service plugin isn't found")
	}

	return p.New(opts...)
}
