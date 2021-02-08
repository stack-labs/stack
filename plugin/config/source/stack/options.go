package stack

import (
	"context"

	"github.com/stack-labs/stack/client"
	"github.com/stack-labs/stack/pkg/config/source"
)

type serviceNameKey struct{}
type pathKey struct{}
type clientKey struct{}

func ServiceName(a string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, serviceNameKey{}, a)
	}
}

// Path sets the key prefix to use
func Path(p string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, pathKey{}, p)
	}
}

func Client(c client.Client) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, clientKey{}, c)
	}
}
