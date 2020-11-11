// Package stack provides a stack rpc resolver which prefixes a namespace
package stack

import (
	"net/http"

	"github.com/stack-labs/stack-rpc/api/resolver"
)

// default resolver for legacy purposes
// it uses proxy routing to resolve names
// /foo becomes namespace.foo
// /v1/foo becomes namespace.v1.foo
type Resolver struct {
	Options resolver.Options
}

func (r *Resolver) Resolve(req *http.Request) (*resolver.Endpoint, error) {
	var name, method string

	switch r.Options.Handler {
	// internal handlers
	case "meta", "api", "rpc", "stack":
		name, method = apiRoute(req.URL.Path)
	default:
		method = req.Method
		name = proxyRoute(req.URL.Path)
	}

	return &resolver.Endpoint{
		Name:   name,
		Method: method,
	}, nil
}

func (r *Resolver) String() string {
	return "stack"
}

// NewResolver creates a new stack resolver
func NewResolver(opts ...resolver.Option) resolver.Resolver {
	return &Resolver{
		Options: resolver.NewOptions(opts...),
	}
}
