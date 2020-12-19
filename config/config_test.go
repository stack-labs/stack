package config

import (
	"fmt"
	"github.com/stack-labs/stack-rpc/pkg/config/source/memory"
	"testing"
)

var (
	ymlFile = []byte(`
stack:
  broker:
    name: http
    address: :8081
  registry:
    interval: 8
`)
)

func TestAutowired(t *testing.T) {
	c := NewConfig(Source(memory.NewSource(memory.WithYAML(ymlFile))))
	if err := c.Init(); err != nil {
		t.Error(fmt.Errorf("Config init error: %s ", err))
	}

	testValue := Value{}
	RegisterOptions(&testValue)

	c.Init()

	if testValue.Stack.Broker.Name != "http" {
		t.Fatalf("broker name should be http, but it's %s", testValue.Stack.Broker.Name)
	}

	if testValue.Stack.Registry.Interval.String() != "8" {
		t.Fatalf("registry interval should be 8, but it's %s", testValue.Stack.Registry.Interval.String())
	}
}
