package util

import (
	"fmt"
	"os"

	"github.com/stack-labs/stack/pkg/cli"

)

type Exec func(*cli.Context, []string) ([]byte, error)

func Print(e Exec) func(ctx *cli.Context) error {
	return func(c *cli.Context) error {
		rsp, err := e(c, c.Args().Slice())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(rsp) > 0 {
			fmt.Printf("%s\n", string(rsp))
		}
		return nil
	}
}
