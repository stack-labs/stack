package cli

import (
	"context"

	"github.com/stack-labs/stack-rpc/cli"
	"github.com/stack-labs/stack-rpc/config/source"
)

type contextKey struct{}

// Context sets the cli context
func Context(c *cli.Context) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, contextKey{}, c)
	}
}
