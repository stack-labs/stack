package config

import (
	"fmt"
	"testing"

	"github.com/stack-labs/stack/pkg/config/source/memory"
)

var (
	ymlFile = []byte(`
stack:
  registry:
    interval: 8
  broker:
    name: http
    address: :8081
`)
)

type testV struct {
	Stack struct {
		Registry struct {
			Interval int `sc:"interval"`
		} `sc:"registry"`
		Broker struct {
			Name    string `sc:"name"`
			Address string `sc:"address"`
		} `sc:"broker"`
	} `sc:"stack"`
}

func TestAutowired(t *testing.T) {
	c := NewConfig(Source(memory.NewSource(memory.WithYAML(ymlFile))))
	if err := c.Init(); err != nil {
		t.Error(fmt.Errorf("Config init error: %s ", err))
	}

	testValue := testV{}
	RegisterOptions(&testValue)

	c.Init()

	if testValue.Stack.Broker.Name != "http" {
		t.Fatalf("broker name should be http, but it's %s", testValue.Stack.Broker.Name)
	}

	if testValue.Stack.Broker.Address != ":8081" {
		t.Fatalf("broker name should be :8081, but it's %s", testValue.Stack.Broker.Address)
	}

	if testValue.Stack.Registry.Interval != 8 {
		t.Fatalf("registry interval should be 8, but it's %d", testValue.Stack.Registry.Interval)
	}
}
