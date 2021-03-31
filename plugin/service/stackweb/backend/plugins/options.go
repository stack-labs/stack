package plugins

import (
	"github.com/stack-labs/stack/service"
)

type Option func(o *Options)

type Options struct {
	Service service.Options
}

func Service(s service.Options) Option {
	return func(o *Options) {
		o.Service = s
	}
}
