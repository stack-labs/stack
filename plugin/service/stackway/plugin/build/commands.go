package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/stack-labs/stack-rpc/pkg/cli"
	goplugin "github.com/stack-labs/stack-rpc/plugin"

	"github.com/stack-labs/stack-rpc-plugins/service/stackway/plugin"
)

func build(ctx *cli.Context) {
	name := ctx.String("name")
	path := ctx.String("path")
	newfn := ctx.String("func")
	typ := ctx.String("type")
	out := ctx.String("output")

	if len(name) == 0 {
		fmt.Println("specify --name of plugin")
		os.Exit(1)
	}

	if len(typ) == 0 {
		fmt.Println("specify --type of plugin")
		os.Exit(1)
	}

	// set the path
	if len(path) == 0 {
		// github.com/stack-rpc/stack-rpc-plugins/broker/rabbitmq
		// github.com/stack-rpc/stack-rpc-plugins/stack/basic_auth
		path = filepath.Join("github.com/stack-labs/stack-rpc-plugins", typ, name)
	}

	// set the newfn
	if len(newfn) == 0 {
		if typ == "stack" {
			newfn = "NewPlugin"
		} else {
			newfn = "New" + strings.Title(typ)
		}
	}

	if len(out) == 0 {
		out = "./"
	}

	// create a .so file
	if !strings.HasSuffix(out, ".so") {
		out = filepath.Join(out, name+".so")
	}

	if err := goplugin.Build(out, &goplugin.Config{
		Name:    name,
		Type:    typ,
		Path:    path,
		NewFunc: newfn,
	}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Plugin %s generated at %s\n", name, out)
}

func pluginCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "build",
			Usage:  "Build a stack plugin",
			Action: build,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "Name of the plugin e.g rabbitmq",
				},
				cli.StringFlag{
					Name:  "type",
					Usage: "Type of the plugin e.g broker",
				},
				cli.StringFlag{
					Name:  "path",
					Usage: "Import path of the plugin",
				},
				cli.StringFlag{
					Name:  "func",
					Usage: "New plugin function creator name e.g NewBroker",
				},
				cli.StringFlag{
					Name:  "output, o",
					Usage: "Output dir or file for the plugin",
				},
			},
		},
	}
}

// Commands returns license commands
func Commands() []cli.Command {
	return []cli.Command{{
		Name:        "plugin",
		Usage:       "Plugin commands",
		Subcommands: pluginCommands(),
	}}
}

// returns a stack plugin which loads plugins
func Flags() plugin.Plugin {
	return plugin.NewPlugin(
		plugin.WithName("plugin"),
		plugin.WithFlag(
			cli.StringSliceFlag{
				Name:   "plugin",
				EnvVar: "STACK_PLUGIN",
				Usage:  "Comma separated list of plugins e.g broker/rabbitmq, registry/etcd, micro/basic_auth, /path/to/plugin.so",
			},
		),
	)
}
