package config

import (
	"github.com/stack-labs/stack-rpc/config/loader"
	"github.com/stack-labs/stack-rpc/config/reader"
	"github.com/stack-labs/stack-rpc/config/source"
)

// Loader sets the loader for manager config
func Loader(l loader.Loader) Option {
	return func(o *Options) {
		o.Loader = l
	}
}

// Source appends a source to list of sources
func Source(s source.Source) Option {
	return func(o *Options) {
		o.Source = append(o.Source, s)
	}
}

// Reader sets the config reader
func Reader(r reader.Reader) Option {
	return func(o *Options) {
		o.Reader = r
	}
}
