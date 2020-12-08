// Package api provides an http-rpc handler which provides the entire http request over rpc
package api

import (
	"net/http"

	goapi "github.com/stack-labs/stack-rpc/api"
	"github.com/stack-labs/stack-rpc/api/handler"
	api "github.com/stack-labs/stack-rpc/api/proto"
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/util/ctx"
	"github.com/stack-labs/stack-rpc/util/errors"
)

type apiHandler struct {
	opts handler.Options
	s    *goapi.Service
}

const (
	Handler = "api"
)

// API handler is the default handler which takes api.Request and returns api.Response
func (a *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request, err := requestToProto(r)
	if err != nil {
		er := errors.InternalServerError("stack.rpc.api", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(er.Error()))
		return
	}

	var service *goapi.Service

	if a.s != nil {
		// we were given the service
		service = a.s
	} else if a.opts.Router != nil {
		// try get service from router
		s, err := a.opts.Router.Route(r)
		if err != nil {
			er := errors.InternalServerError("stack.rpc.api", err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write([]byte(er.Error()))
			return
		}
		service = s
	} else {
		// we have no way of routing the request
		er := errors.InternalServerError("stack.rpc.api", "no route found")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(er.Error()))
		return
	}

	// create request and response
	c := a.opts.Service.Client()
	req := c.NewRequest(service.Name, service.Endpoint.Name, request)
	rsp := &api.Response{}

	// create the context from headers
	cx := ctx.FromRequest(r)
	// create strategy
	so := selector.WithStrategy(strategy(service.Services))

	if err := c.Call(cx, req, rsp, client.WithSelectOption(so)); err != nil {
		w.Header().Set("Content-Type", "application/json")
		ce := errors.Parse(err.Error())
		switch ce.Code {
		case 0:
			w.WriteHeader(500)
		default:
			w.WriteHeader(int(ce.Code))
		}
		w.Write([]byte(ce.Error()))
		return
	} else if rsp.StatusCode == 0 {
		rsp.StatusCode = http.StatusOK
	}

	for _, header := range rsp.GetHeader() {
		for _, val := range header.Values {
			w.Header().Add(header.Key, val)
		}
	}

	if len(w.Header().Get("Content-Type")) == 0 {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(int(rsp.StatusCode))
	w.Write([]byte(rsp.Body))
}

func (a *apiHandler) String() string {
	return "api"
}

func NewHandler(opts ...handler.Option) handler.Handler {
	options := handler.NewOptions(opts...)
	return &apiHandler{
		opts: options,
	}
}

func WithService(s *goapi.Service, opts ...handler.Option) handler.Handler {
	options := handler.NewOptions(opts...)
	return &apiHandler{
		opts: options,
		s:    s,
	}
}
