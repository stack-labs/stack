package web

import (
	"context"
	"net/http"

	"github.com/stack-labs/stack-rpc/auth"
	broker "github.com/stack-labs/stack-rpc/broker/http"
	client "github.com/stack-labs/stack-rpc/client/http"
	selectorR "github.com/stack-labs/stack-rpc/client/selector/registry"
	"github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/registry/mdns"
	"github.com/stack-labs/stack-rpc/server"
	serverH "github.com/stack-labs/stack-rpc/server/http"
	"github.com/stack-labs/stack-rpc/service"
	transportH "github.com/stack-labs/stack-rpc/transport/http"
)

type enableKey struct{}
type addrKey struct{}
type staticDirKey struct{}
type rootPathKey struct{}
type handlersKey struct{}
type serverMuxKey struct{}
type handlerFuncsKey struct{}
type handlerOptsKey struct{}

type HandlerFunc struct {
	Route string
	Func  func(w http.ResponseWriter, r *http.Request)
}

type staticDir struct {
	Route string
	Dir   string
}

func Enable(b bool) service.Option {
	return setOption(enableKey{}, b)
}

func Address(addr string) service.Option {
	return setOption(addrKey{}, addr)
}

func StaticDir(path, dir string) service.Option {
	return setOption(staticDirKey{}, staticDir{
		Route: path,
		Dir:   dir,
	})
}

func RootPath(path string) service.Option {
	return setOption(rootPathKey{}, path)
}

func Handlers(hs ...http.Handler) service.Option {
	return func(o *service.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		v, ok := o.Context.Value(handlersKey{}).([]http.Handler)
		if ok {
			v = append(v, hs...)
		} else {
			v = hs
		}

		o.Context = context.WithValue(o.Context, handlersKey{}, v)
	}
}

func ServerMux(mux *http.ServeMux) service.Option {
	return func(o *service.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, serverMuxKey{}, mux)
	}
}

func HandleFuncs(hs ...HandlerFunc) service.Option {
	return func(o *service.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		v, ok := o.Context.Value(handlerFuncsKey{}).([]HandlerFunc)
		if ok {
			v = append(v, hs...)
		} else {
			v = hs
		}

		o.Context = context.WithValue(o.Context, handlerFuncsKey{}, v)
	}
}

func HandlerOptions(opts ...server.HandlerOption) service.Option {
	return func(o *service.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		v, ok := o.Context.Value(handlerOptsKey{}).([]server.HandlerOption)
		if ok {
			v = append(v, opts...)
		} else {
			v = opts
		}

		o.Context = context.WithValue(o.Context, handlerOptsKey{}, v)
	}
}

func setOption(k, v interface{}) service.Option {
	return func(o *service.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

func newOptions(opts ...service.Option) []service.Option {
	defaultOptions := []service.Option{
		service.Broker(broker.NewBroker()),
		service.Client(client.NewClient()),
		service.Registry(mdns.NewRegistry()),
		service.Server(serverH.NewServer()),
		service.Transport(transportH.NewTransport()),
		service.Selector(selectorR.NewSelector()),
		service.Logger(logger.DefaultLogger),
		service.Config(config.DefaultConfig),
		service.Auth(auth.NoopAuth),
		service.Context(context.Background()),
		service.HandleSignal(true),
	}

	return append(defaultOptions, opts...)
}
