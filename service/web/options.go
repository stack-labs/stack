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
	server "github.com/stack-labs/stack-rpc/server/http"
	"github.com/stack-labs/stack-rpc/service"
	transportH "github.com/stack-labs/stack-rpc/transport/http"
)

type enableKey struct{}
type addrKey struct{}
type staticDirKey struct{}
type rootPathKey struct{}
type handlersKey struct{}
type handlerFuncsKey struct{}

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func Enable(b bool) service.Option {
	return setOption(enableKey{}, b)
}

func Address(addr string) service.Option {
	return setOption(addrKey{}, addr)
}

func StaticDir(dir string) service.Option {
	return setOption(staticDirKey{}, dir)
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

func setOption(k, v interface{}) service.Option {
	return func(o *service.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

func newOptions(opts ...service.Option) service.Options {
	opt := service.Options{
		Broker:    broker.NewBroker(),
		Client:    client.NewClient(),
		Registry:  mdns.NewRegistry(),
		Server:    server.NewServer(),
		Transport: transportH.NewTransport(),
		Selector:  selectorR.NewSelector(),
		Logger:    logger.DefaultLogger,
		Config:    config.DefaultConfig,
		Auth:      auth.NoopAuth,
		Context:   context.Background(),
		Signal:    true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}
