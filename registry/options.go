package registry

import (
	"context"
	"crypto/tls"
	"time"
)

type Options struct {
	Addrs     []string
	Timeout   time.Duration
	Secure    bool
	TLSConfig *tls.Config
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type RegisterOptions struct {
	TTL time.Duration
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
	Domain  string
}

type WatchOptions struct {
	// Specify a service to watch
	// If blank, the watch is for all services
	Service string
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type Option func(*Options)

type DeregisterOptions struct {
	Context context.Context
	// Domain the service was registered in
	Domain string
}

type GetOptions struct {
	Context context.Context
	// Domain to scope the request to
	Domain string
}

type ListOptions struct {
	Context context.Context
	// Domain to scope the request to
	Domain string
}

type RegisterOption func(*RegisterOptions)
type DeregisterOption func(*DeregisterOptions)
type WatchOption func(*WatchOptions)
type GetOption func(*GetOptions)
type ListOption func(*ListOptions)

// Addrs is the registry addresses to use
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

// Secure communication with the registry
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// Specify TLS Config
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

func RegisterTTL(t time.Duration) RegisterOption {
	return func(o *RegisterOptions) {
		o.TTL = t
	}
}

// Watch a service
func WatchService(name string) WatchOption {
	return func(o *WatchOptions) {
		o.Service = name
	}
}
