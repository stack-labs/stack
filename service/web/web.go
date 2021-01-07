package web

import (
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"github.com/stack-labs/stack-rpc"
	broker "github.com/stack-labs/stack-rpc/broker/http"
	cl "github.com/stack-labs/stack-rpc/client"
	client "github.com/stack-labs/stack-rpc/client/http"
	"github.com/stack-labs/stack-rpc/debug/handler"
	ser "github.com/stack-labs/stack-rpc/server"
	server "github.com/stack-labs/stack-rpc/server/http"
	"github.com/stack-labs/stack-rpc/service"
	"github.com/stack-labs/stack-rpc/util/log"
	signalutil "github.com/stack-labs/stack-rpc/util/signal"
)

type webService struct {
	opts service.Options

	once sync.Once
}

func (w *webService) Name() string {
	return w.opts.Server.Options().Name
}

func (w *webService) Init(option ...service.Option) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`hello world`))
	})

	w.Server().NewHandler(mux)

	return nil
}

func (w *webService) Options() service.Options {
	return w.opts
}

func (w *webService) Client() cl.Client {
	return w.opts.Client
}

func (w *webService) Server() ser.Server {
	return w.opts.Server
}

func (w *webService) Run() error {
	// register the debug handler
	w.opts.Server.Handle(
		w.opts.Server.NewHandler(
			handler.NewHandler(w.opts.Client),
			ser.InternalHandler(true),
		),
	)

	// start the profiler
	if w.opts.Profile != nil {
		// to view mutex contention
		runtime.SetMutexProfileFraction(5)
		// to view blocking profile
		runtime.SetBlockProfileRate(1)

		if err := w.opts.Profile.Start(); err != nil {
			return err
		}
		defer w.opts.Profile.Stop()
	}

	log.Infof("Starting [service] %s", w.Name())

	if err := w.Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	if w.opts.Signal {
		signal.Notify(ch, signalutil.Shutdown()...)
	}

	select {
	// wait on kill signal
	case <-ch:
	// wait on context cancel
	case <-w.opts.Context.Done():
	}

	return w.Stop()
}

func (w *webService) Start() error {
	for _, fn := range w.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	if err := w.opts.Server.Start(); err != nil {
		return err
	}

	for _, fn := range w.opts.AfterStart {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (w *webService) Stop() error {
	var gerr error

	for _, fn := range w.opts.BeforeStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	if err := w.opts.Server.Stop(); err != nil {
		return err
	}

	for _, fn := range w.opts.AfterStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	return gerr
}

func (w *webService) String() string {
	return "web"
}

func NewService(opts ...service.Option) service.Service {
	service := new(webService)
	service.opts = newOptions(opts...)
	return service
}

// NewFunction returns a grpc service compatible with stack-rpc.Function
func NewFunction(opts ...service.Option) service.Function {
	c := client.NewClient()
	s := server.NewServer()
	b := broker.NewBroker()

	options := []service.Option{
		stack.Client(c),
		stack.Server(s),
		stack.Broker(b),
		stack.RegisterTTL(time.Minute),
		stack.RegisterInterval(time.Second * 30),
	}

	options = append(options, opts...)

	return stack.NewFunction(options...)
}
