package web

import (
	"github.com/stack-labs/stack-rpc/service"
	"github.com/stack-labs/stack-rpc/service/stack"
	"net/http"
	"path"
	"strings"
)

func NewService(opts ...service.Option) service.Service {
	options := newOptions(opts...)

	options = append(options,
		service.BeforeInit(setHandle),
	)

	return stack.NewService(options...)
}

func setHandle(sOpts *service.Options) error {
	var mux http.Handler
	if sOpts.Context.Value(serverMuxKey{}) != nil {
		if muxTmp, ok := sOpts.Context.Value(serverMuxKey{}).(http.Handler); ok {
			mux = muxTmp
		}
	} else {
		muxTmp := http.NewServeMux()

		rootPath := "/"
		if sOpts.Context.Value(rootPathKey{}) != nil {
			if rootPathTmp, ok := sOpts.Context.Value(rootPathKey{}).(string); ok {
				if !strings.HasPrefix(rootPathTmp, "/") {
					rootPathTmp = "/" + rootPathTmp
				}
				rootPath = rootPathTmp
			}
		}

		// handler funcs
		if sOpts.Context.Value(handlerFuncsKey{}) != nil {
			if handlers, ok := sOpts.Context.Value(handlerFuncsKey{}).([]HandlerFunc); ok {
				for _, handler := range handlers {
					muxTmp.HandleFunc(path.Join(rootPath, handler.Route), handler.Func)
				}
			}
		}

		// static dir
		if sOpts.Context.Value(staticDirKey{}) != nil {
			if sd, ok := sOpts.Context.Value(staticDirKey{}).(staticDir); ok {
				route := path.Join(rootPath, sd.Route)
				muxTmp.Handle(route, http.StripPrefix(route, http.FileServer(http.Dir(sd.Dir))))
			}
		}

		mux = muxTmp
	}

	return sOpts.Server.Handle(sOpts.Server.NewHandler(mux))
}
