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
	Name    string `json:"name"`
	Address string `json:"address"`
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
	Pool    Pool          `json:"pool"`
	Request ClientRequest `json:"request"`
}

type Registry struct {
	Name     string `json:"name"`
	Interval int    `json:"interval"`
	TTL      int    `json:"ttl"`
	Address  string `json:"address"`
}

type Metadata map[string]string

type Server struct {
	Address   string   `json:"address"`
	Advertise string   `json:"advertise"`
	ID        string   `json:"id"`
	Metadata  Metadata `json:"metadata"`
	Name      string   `json:"name"`
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
	Server    Server    `json:"server"`
	Registry  Registry  `json:"registry"`
	Transport Transport `json:"transport"`
	Selector  Selector  `json:"selector"`
	Profile   string    `json:"profile"`
	Runtime   string    `json:"runtime"`
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
		t.Error(err)
	}
	_, err = fh.Write(data)
	if err != nil {
		t.Error(err)
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
		t.Fatal(err)
	}

	conf := Value{
		Stack: Stack{
			Server:    Server{Name: "server-name"},
			Client:    Client{},
			Broker:    Broker{},
			Registry:  Registry{},
			Transport: Transport{},
			Selector:  Selector{},
		},
	}

	conf.Stack.Server.Name = "default-srv-name"
	if err := c.Scan(&conf); err != nil {
		t.Fatal(err)
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
	data1 := []byte(`
---
broker: '1'
transport_address: '1'
`)
	path1 := filepath.Join(os.TempDir(), "file.yaml")
	fh1, err := os.Create(path1)
	if err != nil {
		t.Error(err)
	}
	_, err = fh1.Write(data1)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		fh1.Close()
		os.Remove(path1)
	}()

	// 2 config
	data2 := []byte(`
{ "db" : "mysql"}
`)
	path2 := filepath.Join(os.TempDir(), "file.json")
	fh2, err := os.Create(path2)
	if err != nil {
		t.Error(err)
	}
	_, err = fh2.Write(data2)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		fh2.Close()
		os.Remove(path2)
	}()

	// setup app
	app := cmd.NewCmd().App()
	app.Name = "testcmd"
	app.Flags = cmd.DefaultFlags

	// set args
	os.Args = []string{"run"}
	os.Args = append(os.Args, "--broker", "2")

	c, err := New(path1, app, file.NewSource(file.WithPath(path2)))
	if err != nil {
		t.Fatal(err)
	}

	if c.Get("db").String("default") != "mysql" {
		t.Fatal()
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
		t.Fatal(err)
	}

	if conf.Stack.Server.Name != "default" {
		t.Fatal()
	}

	if conf.Stack.Broker.Name != "http" {
		t.Errorf("broker name [%s] should be 'http'", conf.Stack.Broker.Name)
		t.Fatal()
	}
}
