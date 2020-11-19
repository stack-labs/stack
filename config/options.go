package config

import (
	"github.com/stack-labs/stack-rpc/pkg/cli"
	"github.com/stack-labs/stack-rpc/pkg/config/source"
)

type Options struct {
	FilePath string
	App      *cli.App
	Sources  []source.Source
}

type Option func(o *Options)

func FilePath(f string) Option {
	return func(o *Options) {
		o.FilePath = f
	}
}

func Source(s ...source.Source) Option {
	return func(o *Options) {
		o.Sources = append(o.Sources, s...)
	}
}

func App(a *cli.App) Option {
	return func(o *Options) {
		o.App = a
	}
}
