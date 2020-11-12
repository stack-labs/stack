package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stack-labs/stack-rpc/cmd"
)

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
server_name: '1'
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

	t.Log(c.Config())

	if c.Config().Broker != "2" {
		t.Fatal()
	}
	if c.Config().BrokerAddress != "1" {
		t.Fatal()
	}
	if c.Config().Client != "1" {
		t.Fatal()
	}
	if c.Config().ClientPoolSize != 1 {
		t.Fatal()
	}
	if c.Config().ClientPoolTTL != "1" {
		t.Fatal()
	}
	if c.Config().ClientRequestTimeout != "1" {
		t.Fatal()
	}
	if c.Config().ClientRetries != 1 {
		t.Fatal()
	}
	if c.Config().Profile != "1" {
		t.Fatal()
	}
	if c.Config().RegisterInterval != 1 {
		t.Fatal()
	}
	if c.Config().RegisterTTL != 1 {
		t.Fatal()
	}
	if c.Config().Registry != "1" {
		t.Fatal()
	}
	if c.Config().RegistryAddress != "1" {
		t.Fatal()
	}
	if c.Config().Runtime != "1" {
		t.Fatal()
	}
	if c.Config().Selector != "1" {
		t.Fatal()
	}
	if c.Config().Server != "1" {
		t.Fatal()
	}
	if c.Config().ServerAddress != "1" {
		t.Fatal()
	}
	if c.Config().ServerAdvertise != "1" {
		t.Fatal()
	}
	if c.Config().ServerID != "1" {
		t.Fatal()
	}
	if c.Config().ServerMetadata[0] != "1" {
		t.Fatal()
	}
	if c.Config().ServerName != "1" {
		t.Fatal()
	}
	if c.Config().ServerVersion != "1" {
		t.Fatal()
	}
	if c.Config().Transport != "1" {
		t.Fatal()
	}
	if c.Config().TransportAddress != "1" {
		t.Fatal()
	}
}
