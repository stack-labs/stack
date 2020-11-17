package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stack-labs/stack-rpc/cmd"
	"github.com/stack-labs/stack-rpc/pkg/config/source/file"
)

type Broker struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type Pool struct {
	Size int `json:"size"`
	TTL  int `json:"ttl"`
}

type ClientRequest struct {
	Retries int `json:"retries"`
	Timeout int `json:"timeout"`
}

type Client struct {
	Protocol string        `json:"protocol"`
	Pool     Pool          `json:"pool"`
	Request  ClientRequest `json:"request"`
}

type Registry struct {
	Address  string `json:"address"`
	Interval int    `json:"interval"`
	Name     string `json:"name"`
	TTL      int    `json:"ttl"`
}

type Metadata map[string]string

type Server struct {
	Address   string   `json:"address"`
	Advertise string   `json:"advertise"`
	ID        string   `json:"id"`
	Metadata  Metadata `json:"metadata"`
	Name      string   `json:"name"`
	Protocol  string   `json:"protocol"`
	Version   string   `json:"version"`
}

type Selector struct {
	Name string `json:"name"`
}

type Transport struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Stack struct {
	Broker    Broker    `json:"broker"`
	Client    Client    `json:"client"`
	Profile   string    `json:"profile"`
	Registry  Registry  `json:"registry"`
	Runtime   string    `json:"runtime"`
	Server    Server    `json:"server"`
	Selector  Selector  `json:"selector"`
	Transport Transport `json:"transport"`
}

type Value struct {
	Stack Stack `json:"stack"`
}

func TestStackConfig_Config(t *testing.T) {
	data := []byte(`
stack:
  broker:
    name: http
    address: :8081
  client:
    pool:
      size: 2
      ttl: 200
    request:
      timeout: 300
      retries: 3
  registry:
    name: mdns
    interval: 200
    ttl: 300
    address: 127.0.0.1:6500
  server:
    name:
    address: :8090
    advertise: no-test
    id: test-id
    metadata:
      A: a
      B: b
    version: 1.0.0
  selector:
    name: robin
  transport:
    name: gRPC
    address: :7788
  profile: _1
  runtime:
`)
	path := filepath.Join(os.TempDir(), "file.yaml")
	fh, err := os.Create(path)
	if err != nil {
		t.Error(fmt.Errorf("Config create tmp yml error: %s ", err))
	}
	_, err = fh.Write(data)
	if err != nil {
		t.Error(fmt.Errorf("Config write tmp yml error: %s ", err))
	}
	defer func() {
		fh.Close()
		os.Remove(path)
	}()

	// setup app
	app := cmd.NewCmd().App()
	app.Name = "testcmd"
	app.Flags = cmd.DefaultFlags

	// set args
	os.Args = []string{"run"}
	os.Args = append(os.Args, "--broker", "http", "--broker_address", ":10086")

	c, err := New(path, app)
	if err != nil {
		t.Error(fmt.Errorf("Config new config error: %s ", err))
	}

	conf := Value{
		Stack: Stack{
			Broker:    Broker{},
			Client:    Client{},
			Registry:  Registry{},
			Selector:  Selector{},
			Server:    Server{Name: "server-name"},
			Transport: Transport{},
		},
	}

	conf.Stack.Server.Name = "default-srv-name"
	if err := c.Scan(&conf); err != nil {
		t.Error(fmt.Errorf("Config scan confi error: %s ", err))
	}
	t.Log(conf)

	// test default
	if conf.Stack.Server.Name != "default-srv-name" {
		t.Fatal(fmt.Errorf("server name should be [default-srv-name], not: [%s]", conf.Stack.Server.Name))
	}

	// test the config from cmd
	if conf.Stack.Broker.Address != ":10086" {
		t.Fatal(fmt.Errorf("broker address should be [:10086] which is cmd value, not: [%s]", conf.Stack.Broker.Address))
	}

	// test config deep path
	if conf.Stack.Client.Pool.TTL != 200 {
		t.Fatal(fmt.Errorf("client pool's ttl should be [200], not: [%d]", conf.Stack.Client.Pool.TTL))
	}

	// test config root path
	if conf.Stack.Profile != "_1" {
		t.Fatal(fmt.Errorf("stack profile should be [\"1\"], not: [%s]", conf.Stack.Profile))
	}

	// test map value: the first value
	if conf.Stack.Server.Metadata["A"] != "a" {
		t.Fatal(fmt.Errorf("stack metadata should have [A-a], not: [%s]", conf.Stack.Server.Metadata["A"]))
	}
	// test map value: the extra values
	if conf.Stack.Server.Metadata["B"] != "b" {
		t.Fatal(fmt.Errorf("stack metadata should have [B-b], not: [%s]", conf.Stack.Server.Metadata["B"]))
	}
}

func TestStackConfig_MultiConfig(t *testing.T) {
	// 1 config
	ymlData := []byte(`
stack:
  broker:
    name: http
    address: :8081
  transport:
    name: gRPC
    address: :7788
`)
	ymlPath := filepath.Join(os.TempDir(), "file.yaml")
	ymlFile, err := os.Create(ymlPath)
	if err != nil {
		t.Error(fmt.Errorf("MultiConfig create tmp yml error: %s", err))
	}
	_, err = ymlFile.Write(ymlData)
	if err != nil {
		t.Error(fmt.Errorf("MultiConfig write tmp yml error: %s", err))
	}
	defer func() {
		ymlFile.Close()
		os.Remove(ymlPath)
	}()

	// 2 config
	jsonData := []byte(`
{ "db" : "mysql"}
`)
	jsonPath := filepath.Join(os.TempDir(), "file.json")
	jsonFile, err := os.Create(jsonPath)
	if err != nil {
		t.Error(fmt.Errorf("MultiConfig create tmp json error: %s", err))
	}
	_, err = jsonFile.Write(jsonData)
	if err != nil {
		t.Error(fmt.Errorf("MultiConfig write tmp json error: %s", err))
	}
	defer func() {
		jsonFile.Close()
		os.Remove(jsonPath)
	}()

	// setup app
	app := cmd.NewCmd().App()
	app.Name = "testcmd"
	app.Flags = cmd.DefaultFlags

	// set args
	os.Args = []string{"run"}
	os.Args = append(os.Args, "--broker", "kafka")

	c, err := New(ymlPath, app, file.NewSource(file.WithPath(jsonPath)))
	if err != nil {
		t.Fatal(fmt.Errorf("new config error: %s", err))
	}

	if c.Get("db").String("default") != "mysql" {
		t.Fatal(fmt.Errorf("db setting should be 'mysql', not %s", c.Get("db").String("default")))
	}

	var conf Value
	conf = Value{
		Stack: Stack{
			Server: Server{
				Name: "default",
			},
		},
	}

	if err := c.Scan(&conf); err != nil {
		t.Fatal(fmt.Errorf("MultiConfig scan conf error %s", err))
	}

	if conf.Stack.Server.Name != "default" {
		t.Fatal(fmt.Errorf("broker name [%s] should be 'default'", conf.Stack.Server.Name))
	}

	if conf.Stack.Broker.Name != "kafka" {
		t.Fatal(fmt.Errorf("broker name [%s] should be 'http'", conf.Stack.Broker.Name))
	}
}
