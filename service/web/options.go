package web

import (
	"context"
	"net/http"

	"github.com/stack-labs/stack-rpc"
)

type staticDirKey struct{}
type rootPathKey struct{}
type handlersKey struct{}
type handlerFuncsKey struct{}

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func StaticDir(dir string) stack.Option {
	return setOption(staticDirKey{}, dir)
}

func RootPath(path string) stack.Option {
	return setOption(rootPathKey{}, path)
}

func Handlers(hs ...http.Handler) stack.Option {
	return func(o *stack.Options) {
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

func HandleFuncs(hs ...HandlerFunc) stack.Option {
	return func(o *stack.Options) {
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

func setOption(k, v interface{}) stack.Option {
	return func(o *stack.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
