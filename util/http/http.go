package http

import (
	"net/http"

	"github.com/stack-labs/stack/registry/mdns"

	"github.com/stack-labs/stack/client/selector"
)

func NewRoundTripper(opts ...Option) http.RoundTripper {
	options := Options{
		Registry: mdns.NewRegistry(),
	}
	for _, o := range opts {
		o(&options)
	}

	return &roundTripper{
		rt:   http.DefaultTransport,
		st:   selector.Random(),
		opts: options,
	}
}
