package web

import (
	"github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/plugin"
	"github.com/stack-labs/stack-rpc/service"
)

var options struct {
	Stack struct {
		Service struct {
			Web struct {
				Enable    bool   `sc:"enable"`
				Address   string `sc:"address"`
				RootPath  string `sc:"root-path"`
				StaticDir string `sc:"static-dir"`
			} `sc:"web"`
		} `sc:"service"`
	} `sc:"stack"`
}

type webServicePlugin struct {
}

func (w *webServicePlugin) Name() string {
	return "web"
}

func (w *webServicePlugin) Options() []service.Option {
	var opts []service.Option

	opts = append(opts, Enable(options.Stack.Service.Web.Enable))
	opts = append(opts, Address(options.Stack.Service.Web.Address))
	opts = append(opts, StaticDir(options.Stack.Service.Web.StaticDir))
	opts = append(opts, RootPath(options.Stack.Service.Web.RootPath))

	return opts
}

func (w *webServicePlugin) New(opts ...service.Option) service.Service {
	return NewService(opts...)
}

func init() {
	config.RegisterOptions(&options)
	plugin.ServicePlugins["web"] = &webServicePlugin{}
}
