package grpc

import (
	"time"

	"github.com/stack-labs/stack-rpc"
	"github.com/stack-labs/stack-rpc/broker/http"
	client "github.com/stack-labs/stack-rpc/client/grpc"
	server "github.com/stack-labs/stack-rpc/server/grpc"
)

// NewService returns a grpc service compatible with stack-rpc.Service
func NewService(opts ...stack.Option) stack.Service {
	// our grpc client
	c := client.NewClient()
	// our grpc server
	s := server.NewServer()
	// our grpc broker
	b := http.NewBroker()

	// create options with priority for our opts
	options := []stack.Option{
		stack.Client(c),
		stack.Server(s),
		stack.Broker(b),
	}

	// append passed in opts
	options = append(options, opts...)

	// generate and return a service
	return stack.NewService(options...)
}

// NewFunction returns a grpc service compatible with stack-rpc.Function
func NewFunction(opts ...stack.Option) stack.Function {
	// our grpc client
	c := client.NewClient()
	// our grpc server
	s := server.NewServer()
	// our grpc broker
	b := http.NewBroker()

	// create options with priority for our opts
	options := []stack.Option{
		stack.Client(c),
		stack.Server(s),
		stack.Broker(b),
		stack.RegisterTTL(time.Minute),
		stack.RegisterInterval(time.Second * 30),
	}

	// append passed in opts
	options = append(options, opts...)

	// generate and return a function
	return stack.NewFunction(options...)
}
