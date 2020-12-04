package config

import (
	"github.com/stack-labs/stack-rpc/pkg/config/source"
)

type Options struct {
	Sources []source.Source
	Storage bool
	Watch   bool
}

type Option func(o *Options)

func Source(s ...source.Source) Option {
	return func(o *Options) {
		o.Sources = append(o.Sources, s...)
	}
}

func Storage(s bool) Option {
	return func(o *Options) {
		o.Storage = s
	}
}

func Watch(w bool) Option {
	return func(o *Options) {
		o.Watch = w
	}
}
