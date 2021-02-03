package etcd

import (
	"github.com/stack-labs/stack-rpc/registry"
)

var options struct {
	Stack struct {
		Registry struct {
			Etcd struct {
				AuthCreds struct {
					Username string `sc:"username"`
					Password string `sc:"password"`
				} `sc:"auth-creds"`
			} `sc:"etcd"`
		} `sc:"registry"`
	} `sc:"stack"`
}

type etcdRegistryPlugin struct{}

func (c *etcdRegistryPlugin) Name() string {
	return "etcd"
}

func (c *etcdRegistryPlugin) Options() []registry.Option {
	var opts []registry.Option
	ec := options.Stack.Registry.Etcd

	if len(ec.AuthCreds.Username) > 0 {
		opts = append(opts, Auth(ec.AuthCreds.Username, ec.AuthCreds.Password))
	}

	return opts
}

func (c *etcdRegistryPlugin) New(opts ...registry.Option) registry.Registry {
	return NewRegistry(opts...)
}
