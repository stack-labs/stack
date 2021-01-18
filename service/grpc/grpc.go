package grpc

import (
	"github.com/stack-labs/stack-rpc/broker/http"
	client "github.com/stack-labs/stack-rpc/client/grpc"
	server "github.com/stack-labs/stack-rpc/server/grpc"
	"github.com/stack-labs/stack-rpc/service"
)

// NewOptions returns a grpc service options compatible with stack-rpc.Service
func NewOptions(opts ...service.Option) []service.Option {
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

	return append(options, opts...)
}
