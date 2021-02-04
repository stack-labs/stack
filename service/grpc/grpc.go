package grpc

import (
	"github.com/stack-labs/stack/broker/http"
	client "github.com/stack-labs/stack/client/grpc"
	server "github.com/stack-labs/stack/server/grpc"
	"github.com/stack-labs/stack/service"
)

// NewOptions returns a grpc service options compatible with stack.Service
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
