package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/pkg/config/source"
	cliSource "github.com/stack-labs/stack-rpc/pkg/config/source/cli"
	"github.com/stack-labs/stack-rpc/pkg/config/source/file"
	"github.com/stack-labs/stack-rpc/pkg/config/source/memory"
)

var (
	ymlFile = []byte(`
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
      - A=a
      - B=b
    version: 1.0.0
  selector:
    name: robin
  transport:
    name: gRPC
    address: :7788
  profile: _1
  runtime:
`)
	conf = config.Value{
		Stack: config.Stack{
			Broker:    config.Broker{},
			Client:    config.Client{},
			Registry:  config.Registry{},
			Selector:  config.Selector{},
			Server:    config.Server{},
			Transport: config.Transport{},
		},
	}
)

func TestStackConfig_File(t *testing.T) {
	path := filepath.Join(os.TempDir(), "file.yaml")
	fh, err := os.Create(path)
	if err != nil {
		t.Error(fmt.Errorf("Config create tmp yml error: %s ", err))
	}
	_, err = fh.Write(ymlFile)
	if err != nil {
		t.Error(fmt.Errorf("Config write tmp yml error: %s ", err))
	}
	defer func() {
		fh.Close()
		os.Remove(path)
	}()

	c := config.NewConfig(config.Source(file.NewSource(file.WithPath(path))))
	if err = c.Init(); err != nil {
		t.Error(fmt.Errorf("Config init error: %s ", err))
	}

	if err := c.Scan(&conf); err != nil {
		t.Error(fmt.Errorf("Config scan confi error: %s ", err))
	}
	t.Log(string(c.Bytes()))
	t.Log(conf)

	if conf.Stack.Server.Address != ":8090" {
		t.Fatal(fmt.Errorf("server address should be [:8090], not: [%s]", conf.Stack.Server.Address))
	}
}

func TestStackConfig_Config(t *testing.T) {
	path := filepath.Join(os.TempDir(), "file.yaml")
	fh, err := os.Create(path)
	if err != nil {
		t.Error(fmt.Errorf("Config create tmp yml error: %s ", err))
	}
	_, err = fh.Write(ymlFile)
	if err != nil {
		t.Error(fmt.Errorf("Config write tmp yml error: %s ", err))
	}
	defer func() {
		fh.Close()
		os.Remove(path)
	}()

	// setup app
	app := NewCmd().App()
	app.Name = "testcmd"
	app.Flags = DefaultFlags

	// set args
	os.Args = []string{"run"}
	// string arg
	os.Args = append(os.Args, "--broker", "http", "--broker_address", ":10086")
	// int arg
	os.Args = append(os.Args, "--client_pool_ttl", "100")
	// map
	os.Args = append(os.Args, "--server_metadata", "C=c")
	os.Args = append(os.Args, "--server_metadata", "D=d")

	conf.Stack.Server.Name = "default-srv-name"
	defaultBytes, _ := json.Marshal(conf)

	sources := []source.Source{
		memory.NewSource(memory.WithJSON(defaultBytes)),
		file.NewSource(file.WithPath(path)),
		cliSource.NewSource(app, cliSource.Context(app.Context())),
	}

	c := config.NewConfig(config.Source(sources...))
	if err = c.Init(); err != nil {
		t.Error(fmt.Errorf("Config init error: %s ", err))
	}

	if err := c.Scan(&conf); err != nil {
		t.Error(fmt.Errorf("Config scan confi error: %s ", err))
	}
	t.Log(string(c.Bytes()))
	t.Log(conf)

	// test default
	if conf.Stack.Server.Name != "default-srv-name" {
		t.Fatal(fmt.Errorf("server name should be [default-srv-name], not: [%s]", conf.Stack.Server.Name))
	}
	if c.Get("stack", "server", "name").String("") != "default-srv-name" {
		t.Fatal(fmt.Errorf("server name in [c] should be [default-srv-name], not: [%s]", c.Get("stack", "server", "name").String("")))
	}

	if conf.Stack.Server.ID != "test-id" {
		t.Fatal(fmt.Errorf("server id should be [test-id] which is cmd value, not: [%s]", conf.Stack.Server.ID))
	}

	// test the config from cmd
	if conf.Stack.Broker.Address != ":10086" {
		t.Fatal(fmt.Errorf("broker address should be [:10086] which is cmd value, not: [%s]", conf.Stack.Broker.Address))
	}

	if conf.Stack.Client.Pool.TTL != "100" {
		t.Fatal(fmt.Errorf("client pool ttl should be [100] which is cmd value, not: [%s]", conf.Stack.Client.Pool.TTL))
	}

	if conf.Stack.Registry.TTL != "300" {
		t.Fatal(fmt.Errorf("client registry ttl should be [300] which is cmd value, not: [%s]", conf.Stack.Registry.TTL))
	}

	// test config root path
	if conf.Stack.Profile != "_1" {
		t.Fatal(fmt.Errorf("stack profile should be [\"_1\"], not: [%s]", conf.Stack.Profile))
	}

	// test map value: the extra values
	if conf.Stack.Server.Metadata.Value("C") != "c" {
		t.Fatal(fmt.Errorf("stack metadata should have [C-c], not: [%s]", conf.Stack.Server.Metadata.Value("C")))
	}
	// test map value: the cmd value
	if conf.Stack.Server.Metadata.Value("D") != "d" {
		t.Fatal(fmt.Errorf("stack metadata should have [D-d], not: [%s]", conf.Stack.Server.Metadata.Value("D")))
	}
}

func TestStackConfig_MultiConfig(t *testing.T) {
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
	app := NewCmd().App()
	app.Name = "testcmd"
	app.Flags = DefaultFlags

	// set args
	os.Args = []string{"run"}
	os.Args = append(os.Args, "--broker", "kafka")

	sources := []source.Source{
		file.NewSource(file.WithPath(ymlPath)),
		file.NewSource(file.WithPath(jsonPath)),
		cliSource.NewSource(app, cliSource.Context(app.Context())),
	}

	c := config.NewConfig(config.Source(sources...))
	if err = c.Init(); err != nil {
		t.Error(fmt.Errorf("Config init error: %s ", err))
	}

	if err = c.Init(); err != nil {
		t.Fatal(fmt.Errorf("Config init error: %s ", err))
	}

	if c.Get("db").String("default") != "mysql" {
		t.Fatal(fmt.Errorf("db setting should be 'mysql', not %s", c.Get("db").String("default")))
	}

	var conf config.Value
	conf = config.Value{
		Stack: config.Stack{
			Server: config.Server{
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
