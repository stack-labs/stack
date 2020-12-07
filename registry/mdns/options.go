// Package mdns provides a multicast dns registry
package mdns

import (
	"context"

	"github.com/stack-labs/stack-rpc/registry"
)

// Domain sets the mdnsDomain
func Domain(d string) registry.Option {
	return func(o *registry.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "mdns.domain", d)
	}
}
