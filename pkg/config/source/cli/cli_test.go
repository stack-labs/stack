package cli

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stack-labs/stack-rpc/cli"
	"github.com/stack-labs/stack-rpc/cmd"
	"github.com/stack-labs/stack-rpc/pkg/config/source"
)

func test(t *testing.T, withContext bool) {
	var src source.Source

	// setup app
	app := cmd.NewCmd().App()
	app.Name = "testapp"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "db-host"},
	}

	// with context
	if withContext {
		// set action
		app.Action = func(c *cli.Context) {
			src = WithContext(c)
		}

		// run app
		if err := app.Run([]string{"run", "-db-host", "localhost"}); err != nil {
			t.Error(err)
		}
		// no context
	} else {
		// set args
		os.Args = []string{"run", "-db-host", "localhost"}
		src = NewSource(app)
	}

	// test config
	c, err := src.Read()
	if err != nil {
		t.Error(err)
	}
	if len(c.Data) == 0 {
		t.Fatal()
	}

	t.Log(string(c.Data))

	var actual map[string]interface{}
	if err := json.Unmarshal(c.Data, &actual); err != nil {
		t.Error(err)
	}

	actualDB := actual["db"].(map[string]interface{})
	if actualDB["host"] != "localhost" {
		t.Errorf("expected localhost, got %v", actualDB["name"])
	}

}

func TestCliSource(t *testing.T) {
	// without context
	test(t, false)
}

func TestCliSourceWithContext(t *testing.T) {
	// with context
	test(t, true)
}
