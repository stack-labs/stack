package web

import (
	"net/http"
	"path"
	"strings"

	"github.com/stack-labs/stack-rpc/service"
	"github.com/stack-labs/stack-rpc/util/log"
)

func NewOptions(opts ...service.Option) []service.Option {
	options := newOptions(opts...)

	options = append(options,
		service.BeforeInit(setHandle),
	)

	return options
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
					log.Debugf("handler url, root: [%s], route: [%s]", rootPath, handler.Route)
					route := path.Join(rootPath, handler.Route)
					if strings.HasSuffix(handler.Route, "/") {
						route += "/"
					}
					muxTmp.HandleFunc(route, handler.Func)
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
