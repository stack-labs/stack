package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stack-labs/stack-rpc/cmd"
	"github.com/stack-labs/stack-rpc/pkg/cli"
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

	if actual["db-host"] != "localhost" {
		t.Errorf("expected localhost, got %v", actual["name"])
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

type config struct {
	Broker               string   `json:"broker"`
	BrokerAddress        string   `json:"broker_address"`
	Client               string   `json:"client"`
	ClientPoolSize       int      `json:"client_pool_size"`
	ClientPoolTTL        string   `json:"client_pool_ttl"`
	ClientRequestTimeout string   `json:"client_request_timeout"`
	ClientRetries        int      `json:"client_retries"`
	Profile              string   `json:"profile"`
	RegisterInterval     int      `json:"register_interval"`
	RegisterTTL          int      `json:"register_ttl"`
	Registry             string   `json:"registry"`
	RegistryAddress      string   `json:"registry_address"`
	Selector             string   `json:"selector"`
	Server               string   `json:"server"`
	ServerAddress        string   `json:"server_address"`
	ServerAdvertise      string   `json:"server_advertise"`
	ServerID             string   `json:"server_id"`
	ServerMetadata       []string `json:"server_metadata"`
	ServerName           string   `json:"server_name"`
	ServerVersion        string   `json:"server_version"`
	Transport            string   `json:"transport"`
	TransportAddress     string   `json:"transport_address"`
}

func TestCliSource_cmd(t *testing.T) {
	// setup app
	app := cmd.NewCmd().App()
	app.Name = "testcmd"
	app.Flags = cmd.DefaultFlags

	// set args
	os.Args = []string{"run"}
	for _, v := range cmd.DefaultFlags {
		os.Args = append(os.Args, fmt.Sprintf("--%s", v.GetName()), "1")
	}
	src := NewSource(app)

	// test config
	c, err := src.Read()
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Data) == 0 {
		t.Fatal()
	}

	t.Log(string(c.Data))

	var conf config
	if err := json.Unmarshal(c.Data, &conf); err != nil {
		t.Error(err)
	}

	if conf.Broker != "1" {
		t.Fatal()
	}
	if conf.BrokerAddress != "1" {
		t.Fatal()
	}
	if conf.Client != "1" {
		t.Fatal()
	}
	if conf.ClientPoolSize != 1 {
		t.Fatal()
	}
	if conf.ClientPoolTTL != "1" {
		t.Fatal()
	}
	if conf.ClientRequestTimeout != "1" {
		t.Fatal()
	}
	if conf.ClientRetries != 1 {
		t.Fatal()
	}
	if conf.Profile != "1" {
		t.Fatal()
	}
	if conf.RegisterInterval != 1 {
		t.Fatal()
	}
	if conf.RegisterTTL != 1 {
		t.Fatal()
	}
	if conf.Registry != "1" {
		t.Fatal()
	}
	if conf.RegistryAddress != "1" {
		t.Fatal()
	}
	if conf.Selector != "1" {
		t.Fatal()
	}
	if conf.Server != "1" {
		t.Fatal()
	}
	if conf.ServerAddress != "1" {
		t.Fatal()
	}
	if conf.ServerAdvertise != "1" {
		t.Fatal()
	}
	if conf.ServerID != "1" {
		t.Fatal()
	}
	if conf.ServerMetadata[0] != "1" {
		t.Fatal()
	}
	if conf.ServerName != "1" {
		t.Fatal()
	}
	if conf.ServerVersion != "1" {
		t.Fatal()
	}
	if conf.Transport != "1" {
		t.Fatal()
	}
	if conf.TransportAddress != "1" {
		t.Fatal()
	}
}
