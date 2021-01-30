package consul

import (
	"time"

	"github.com/stack-labs/stack-rpc/registry"
)

var options struct {
	Stack struct {
		Registry struct {
			Consul struct {
				Connect    bool `sc:"connect"`
				AllowStale bool `sc:"allow-stale"`
				TCPCheck   int  `sc:"tcp-check"`
			} `sc:"consul"`
		} `sc:"registry"`
	} `sc:"stack"`
}

type consulRegistryPlugin struct{}

func (c *consulRegistryPlugin) Name() string {
	return "consul"
}

func (c *consulRegistryPlugin) Options() []registry.Option {
	var opts []registry.Option
	cc := options.Stack.Registry.Consul

	if cc.Connect {
		opts = append(opts, Connect())
	}

	opts = append(opts, AllowStale(cc.AllowStale))

	if ttl := time.Duration(cc.TCPCheck); ttl >= 0 {
		opts = append(opts, TCPCheck(ttl*time.Second))
	}

	return opts
}

func (c *consulRegistryPlugin) New(opts ...registry.Option) registry.Registry {
	return NewRegistry(opts...)
}
