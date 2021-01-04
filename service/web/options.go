package web

import (
	"context"

	"github.com/stack-labs/stack-rpc"
)

type staticDirKey struct{}
type rootPathKey struct{}

func StaticDir(dir string) stack.Option {
	return SetOption(staticDirKey{}, dir)
}

func RootPathKey(path string) stack.Option {
	return SetOption(rootPathKey{}, path)
}

func SetOption(k, v interface{}) stack.Option {
	return func(o *stack.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
