package etcd

import (
	"testing"
	"time"

	"github.com/stack-labs/stack-rpc/registry"
)

var testData = map[string][]*registry.Service{
	"foo": {
		{
			Name:    "foo",
			Version: "1.0.0",
			Nodes: []*registry.Node{
				{
					Id:      "foo-1.0.0-123",
					Address: "localhost:9999",
				},
				{
					Id:      "foo-1.0.0-321",
					Address: "localhost:9999",
				},
			},
		},
		{
			Name:    "foo",
			Version: "1.0.1",
			Nodes: []*registry.Node{
				{
					Id:      "foo-1.0.1-321",
					Address: "localhost:6666",
				},
			},
		},
		{
			Name:    "foo",
			Version: "1.0.3",
			Nodes: []*registry.Node{
				{
					Id:      "foo-1.0.3-345",
					Address: "localhost:8888",
				},
			},
		},
	},
	"bar": {
		{
			Name:    "bar",
			Version: "default",
			Nodes: []*registry.Node{
				{
					Id:      "bar-1.0.0-123",
					Address: "localhost:9999",
				},
				{
					Id:      "bar-1.0.0-321",
					Address: "localhost:9999",
				},
			},
		},
		{
			Name:    "bar",
			Version: "latest",
			Nodes: []*registry.Node{
				{
					Id:      "bar-1.0.1-321",
					Address: "localhost:6666",
				},
			},
		},
	},
}

func TestEtcdRegistryTTL(t *testing.T) {
	m := NewRegistry()

	for _, v := range testData {
		for _, service := range v {
			if err := m.Register(service, registry.RegisterTTL(time.Second*2)); err != nil {
				t.Fatal(err)
			}
		}
	}

	// sleep enough seconds
	time.Sleep(time.Second * 10)

	for name := range testData {
		svcs, err := m.GetService(name)
		if err != nil && err != registry.ErrNotFound {
			t.Fatal(err)
		}

		for _, svc := range svcs {
			if len(svc.Nodes) > 0 {
				t.Fatalf("Service %q still has nodes registered", svc.Nodes[0].Address)
			}
		}
	}
}
