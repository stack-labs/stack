package stack

import (
	"github.com/stack-labs/stack-rpc/service"
	"github.com/stack-labs/stack-rpc/service/web"
)

type webServicePlugin struct{}

func (w *webServicePlugin) Name() string {
	return "web"
}

func (w *webServicePlugin) Options() []service.Option {
	var opts []service.Option

	opts = append(opts, web.Enable(options.Stack.Service.Web.Enable))
	opts = append(opts, web.Address(options.Stack.Service.Web.Address))
	opts = append(opts, web.StaticDir("", options.Stack.Service.Web.StaticDir))
	opts = append(opts, web.RootPath(options.Stack.Service.Web.RootPath))

	return opts
}

func (w *webServicePlugin) New(opts ...service.Option) service.Service {
	return web.NewService(opts...)
}
