// Package registry is an interface for service discovery
package registry

import (
	"context"
	"errors"
)

var (
	DefaultRegistry = NewRegistry()

	// Not found error when GetService is called
	ErrNotFound = errors.New("service not found")
	// Watcher stopped error when watcher is stopped
	ErrWatcherStopped = errors.New("watcher stopped")
)

// The registry provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}
type Registry interface {
	Init(...Option) error
	Options() Options
	Register(*Service, ...RegisterOption) error
	Deregister(*Service) error
	GetService(string) ([]*Service, error)
	ListServices() ([]*Service, error)
	Watch(...WatchOption) (Watcher, error)
	String() string
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

// Register a service node. Additionally supply options such as TTL.
func Register(s *Service, opts ...RegisterOption) error {
	return DefaultRegistry.Register(s, opts...)
}

// Deregister a service node
func Deregister(s *Service) error {
	return DefaultRegistry.Deregister(s)
}

// Retrieve a service. A slice is returned since we separate Name/Version.
func GetService(name string) ([]*Service, error) {
	return DefaultRegistry.GetService(name)
}

// List the services. Only returns service names
func ListServices() ([]*Service, error) {
	return DefaultRegistry.ListServices()
}

// Watch returns a watcher which allows you to track updates to the registry.
func Watch(opts ...WatchOption) (Watcher, error) {
	return DefaultRegistry.Watch(opts...)
}

func String() string {
	return DefaultRegistry.String()
}
