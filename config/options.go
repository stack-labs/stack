package config

import (
	"context"

	"github.com/stack-labs/stack/pkg/config/source"
)

type Options struct {
	Sources []source.Source
	Storage bool
	Watch   bool
	// HierarchyMerge merges the query args to one
	// eg. Get("a","b","c") can be used as Get("a.b.c")
	// the default is false
	HierarchyMerge bool

	Context context.Context
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

func HierarchyMerge(h bool) Option {
	return func(o *Options) {
		o.HierarchyMerge = h
	}
}

func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}
