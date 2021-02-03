package apollo

import (
	"context"

	"github.com/stack-labs/stack-rpc/pkg/config/source"
)

type appIDKey struct{}
type clusterKey struct{}
type addrKey struct{}
type namespacesKey struct{}
type isBackupConfigKey struct{}
type secretKey struct{}

func AppID(id string) source.Option {
	return func(o *source.Options) {
		prepareCtx(o)
		o.Context = context.WithValue(o.Context, appIDKey{}, id)
	}
}

func Cluster(cluster string) source.Option {
	return func(o *source.Options) {
		prepareCtx(o)
		o.Context = context.WithValue(o.Context, clusterKey{}, cluster)
	}
}

func Addr(addr string) source.Option {
	return func(o *source.Options) {
		prepareCtx(o)
		o.Context = context.WithValue(o.Context, addrKey{}, addr)
	}
}

func Namespaces(namespaces string) source.Option {
	return func(o *source.Options) {
		prepareCtx(o)
		o.Context = context.WithValue(o.Context, namespacesKey{}, namespaces)
	}
}

func IsBackupConfig(isBackupConfig bool) source.Option {
	return func(o *source.Options) {
		prepareCtx(o)
		o.Context = context.WithValue(o.Context, isBackupConfigKey{}, isBackupConfig)
	}
}

func Secret(secret string) source.Option {
	return func(o *source.Options) {
		prepareCtx(o)
		o.Context = context.WithValue(o.Context, secretKey{}, secret)
	}
}

func prepareCtx(o *source.Options) {
	if o.Context == nil {
		o.Context = context.Background()
	}
}
