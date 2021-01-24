package main

import (
	"os"

	"github.com/stack-labs/stack-rpc/pkg/cli"
	"github.com/stack-labs/stack-rpc/util/stackctl/new"
	"github.com/stack-labs/stack-rpc/util/stackctl/service"
)

func main() {
	app := cli.NewApp()
	app.Name = "stackctl"

	app.Commands = append(app.Commands, new.Commands()...)
	app.Commands = append(app.Commands, service.Commands()...)

	app.Run(os.Args)
}
