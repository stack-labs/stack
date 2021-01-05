package web

import (
	"github.com/stack-labs/stack-rpc/util/wrapper"
	"net/http"
	"time"

	"github.com/stack-labs/stack-rpc"
	broker "github.com/stack-labs/stack-rpc/broker/http"
	cl "github.com/stack-labs/stack-rpc/client"
	client "github.com/stack-labs/stack-rpc/client/http"
	ser "github.com/stack-labs/stack-rpc/server"
	server "github.com/stack-labs/stack-rpc/server/http"
)

type webService struct {
}

func (w *webService) Name() string {
	return w.Options().Server.Options().Name
}

func (w *webService) Init(option ...stack.Option) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`hello world`))
	})

	w.Server().NewHandler(mux)

	return nil
}

func (w *webService) Options() stack.Options {
	return w.Options()
}

func (w *webService) Client() cl.Client {
	return w.Client()
}

func (w *webService) Server() ser.Server {
	return w.Server()
}

func (w *webService) Run() error {
	panic("implement me")
}

func (w *webService) String() string {
	return "web"
}

func NewService(opts ...stack.Option) stack.Service {
	c := client.NewClient()
	s := server.NewServer()
	b := broker.NewBroker()

	// wrap client to inject From-Service header on any calls
	c = wrapper.FromService(serviceName, options.Client)
	c = wrapper.TraceCall(serviceName, trace.DefaultTracer, options.Client)
	c = wrapper.CacheClient(cacheFn, options.Client)
	c = wrapper.AuthClient(authFn, options.Client)

	options := []stack.Option{
		stack.Client(c),
		stack.Server(s),
		stack.Broker(b),
	}

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
