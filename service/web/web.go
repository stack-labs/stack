package web

import (
	"time"

	"github.com/stack-labs/stack-rpc"
	broker "github.com/stack-labs/stack-rpc/broker/http"
	client "github.com/stack-labs/stack-rpc/client/http"
	server "github.com/stack-labs/stack-rpc/server/http"
)

func NewService(opts ...stack.Option) stack.Service {
	c := client.NewClient()
	s := server.NewServer()
	b := broker.NewBroker()

	options := []stack.Option{
		stack.Client(c),
		stack.Server(s),
		stack.Broker(b),
	}

	s.Handle()

	options = append(options, opts...)

	return stack.NewService(options...)
}

// NewFunction returns a grpc service compatible with stack-rpc.Function
func NewFunction(opts ...stack.Option) stack.Function {
	c := client.NewClient()
	s := server.NewServer()
	b := broker.NewBroker()

	options := []stack.Option{
		stack.Client(c),
		stack.Server(s),
		stack.Broker(b),
		stack.RegisterTTL(time.Minute),
		stack.RegisterInterval(time.Second * 30),
	}

	options = append(options, opts...)

	return stack.NewFunction(options...)
}
