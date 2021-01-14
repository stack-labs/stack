package grpc

import (
	"github.com/stack-labs/stack-rpc/broker/http"
	client "github.com/stack-labs/stack-rpc/client/grpc"
	server "github.com/stack-labs/stack-rpc/server/grpc"
	"github.com/stack-labs/stack-rpc/service"
	"github.com/stack-labs/stack-rpc/service/stack"
)

// NewService returns a grpc service compatible with stack-rpc.Service
func NewService(opts ...service.Option) service.Service {
	// our grpc client
	c := client.NewClient()
	// our grpc server
	s := server.NewServer()
	// our grpc broker
	b := http.NewBroker()

	// create options with priority for our opts
	options := []service.Option{
		service.Client(c),
		service.Server(s),
		service.Broker(b),
	}

	// append passed in opts
	options = append(options, opts...)

	// use stack service for current
	return stack.NewService(options...)
}
