package stack

import (
	"context"

	"github.com/stack-labs/stack-rpc/auth"
	"github.com/stack-labs/stack-rpc/cmd"
	"github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/plugin"
	"github.com/stack-labs/stack-rpc/service"
	"github.com/stack-labs/stack-rpc/util/log"
)

func newOptions(opts ...service.Option) service.Options {
	options := service.Options{
		Cmd: cmd.NewCmd(),
	}

	defaultOpts := []service.Option{
		service.Context(context.Background()),
		// default use stack rpc
		service.RPC("stack"),
		service.BeforeInit(func(sOpts *service.Options) error {
			// cmd helps stack parse command options and reset the options that should work.
			if err := sOpts.Cmd.Init(
				// todo config passed is not cool
				cmd.Config(&sOpts.Config),
				cmd.ServiceOptions(sOpts),
			); err != nil {
				log.Errorf("cmd init error: %s", err)
				return err
			}
			return nil
		}),
		// set the default components
		service.Broker(plugin.BrokerPlugins["http"].New()),
		service.Client(plugin.ClientPlugins["mucp"].New()),
		service.Server(plugin.ServerPlugins["mucp"].New()),
		service.Registry(plugin.RegistryPlugins["mdns"].New()),
		service.Transport(plugin.TransportPlugins["http"].New()),
		service.Selector(plugin.SelectorPlugins["cache"].New()),
		service.Logger(plugin.LoggerPlugins["console"].New()),
		service.Config(config.DefaultConfig),
		service.Auth(auth.NoopAuth),
		service.HandleSignal(true),
	}

	defaultOpts = append(defaultOpts, opts...)

	for _, o := range defaultOpts {
		o(&options)
	}

	return options
}
