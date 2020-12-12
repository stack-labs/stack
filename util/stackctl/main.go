package main

import (
	"os"

	"github.com/stack-labs/stack-rpc/pkg/cli"
	"github.com/stack-labs/stack-rpc/util/stackctl/new"
)

func main() {
	app := cli.NewApp()
	app.Name = "stackctl"

	app.Commands = append(app.Commands, new.Commands()...)

	app.Run(os.Args)
}
