package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stack-labs/stack-rpc/pkg/config/source/file"

	"github.com/stack-labs/stack-rpc/cmd"
)

type Value struct {
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
	Runtime              string   `json:"runtime"`
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

func TestStackConfig_Config(t *testing.T) {
	data := []byte(`
---
broker: '1'
broker_address: '1'
client: '1'
client_pool_size: 1
client_pool_ttl: '1'
client_request_timeout: '1'
client_retries: 1
profile: '1'
register_interval: 1
register_ttl: 1
registry: '1'
registry_address: '1'
runtime: '1'
selector: '1'
server: '1'
server_address: '1'
server_advertise: '1'
server_id: '1'
server_metadata:
- '1'
server_version: '1'
transport: '1'
transport_address: '1'
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
	os.Args = append(os.Args, "--broker", "2")

	c, err := New(path, app)
	if err != nil {
		t.Fatal(err)
	}

	var conf Value
	conf.ServerName = "default"
	if err := c.Scan(&conf); err != nil {
		t.Fatal(err)
	}
	t.Log(conf)

	if conf.ServerName != "default" {
		t.Fatal()
	}
	if conf.Broker != "2" {
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
	if conf.Runtime != "1" {
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
	conf.ServerName = "default"
	if err := c.Scan(&conf); err != nil {
		t.Fatal(err)
	}
	if conf.ServerName != "default" {
		t.Fatal()
	}
	if conf.Broker != "2" {
		t.Fatal()
	}

}
