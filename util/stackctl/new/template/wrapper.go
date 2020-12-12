package template

var (
	WrapperAPI = `package client

import (
	"context"

	"github.com/stack-labs/stack-rpc"
	"github.com/stack-labs/stack-rpc/server"

	// FIXME: replace with service proto path
	{{.Alias}} "path/to/service/proto/{{.Alias}}"
)

type {{.Alias}}Key struct {}

// FromContext retrieves the client from the Context
func {{title .Alias}}FromContext(ctx context.Context) ({{.Alias}}.{{title .Alias}}Service, bool) {
	c, ok := ctx.Value({{.Alias}}Key{}).({{.Alias}}.{{title .Alias}}Service)
	return c, ok
}

// Client returns a wrapper for the {{title .Alias}}Client
func {{title .Alias}}Wrapper(service stack.Service) server.HandlerWrapper {
	// FIXME: replace "stack.rpc.service.{{.Alias}}" with service name
	client := {{.Alias}}.New{{title .Alias}}Service("stack.rpc.service.{{.Alias}}", service.Client())

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, {{.Alias}}Key{}, client)
			return fn(ctx, req, rsp)
		}
	}
}
`
)
