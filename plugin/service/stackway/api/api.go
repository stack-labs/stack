package api

import (
	"github.com/stack-labs/stack-rpc"
	"github.com/stack-labs/stack-rpc/pkg/cli"
	"github.com/stack-labs/stack-rpc/service"
)

// api stackway options
func Options() (options []service.Option) {
	flags := []cli.Flag{
		cli.StringFlag{
			Name:   "stackway_name",
			Usage:  "Stackway name",
			EnvVar: "STACK_STACKWAY_NAME",
			Alias:  "stack_stackway_name",
		},
		cli.StringFlag{
			Name:   "stackway_address",
			Usage:  "Set the stackway address e.g 0.0.0.0:8080",
			EnvVar: "STACK_STACKWAY_ADDRESS",
			Alias:  "stack_stackway_address",
		},
		cli.StringFlag{
			Name:   "stackway_handler",
			Usage:  "Specify the request handler to be used for mapping HTTP requests to services; {api, event, http, rpc}",
			EnvVar: "STACK_STACKWAY_HANDLER",
			Alias:  "stack_stackway_handler",
		},
		cli.StringFlag{
			Name:   "stackway_namespace",
			Usage:  "Set the namespace used by the stackway e.g. stack.rpc.api",
			EnvVar: "STACK_STACKWAY_NAMESPACE",
			Alias:  "stack_stackway_namespace",
		},
		cli.StringFlag{
			Name:   "stackway_resolver",
			Usage:  "Set the hostname resolver used by the stackway {host, path, grpc}",
			EnvVar: "STACK_STACKWAY_RESOLVER",
			Alias:  "stack_stackway_resolver",
		},
		cli.BoolFlag{
			Name:   "stackway_enable_rpc",
			Usage:  "Enable call the backend directly via /rpc",
			EnvVar: "STACK_STACKWAY_ENABLE_RPC",
			Alias:  "stack_stackway_enable_rpc",
		},
	}

	options = append(options, stack.Flags(flags...))

	return
}
