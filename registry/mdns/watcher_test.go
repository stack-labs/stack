package mdns

import (
	"os"
	"testing"
	"time"

	"github.com/stack-labs/stack-rpc/registry"
)

func TestWatcher(t *testing.T) {
	testData := []*registry.Service{
		{
			Name:    "test1",
			Version: "1.0.1",
			Nodes: []*registry.Node{
				{
					Id:      "test1-1",
					Address: "10.0.0.1:10001",
					Metadata: map[string]string{
						"foo": "bar",
					},
				},
			},
		},
		{
			Name:    "test2",
			Version: "1.0.2",
			Nodes: []*registry.Node{
				{
					Id:      "test2-1",
					Address: "10.0.0.2:10002",
					Metadata: map[string]string{
						"foo2": "bar2",
					},
				},
			},
		},
		{
			Name:    "test3",
			Version: "1.0.3",
			Nodes: []*registry.Node{
				{
					Id:      "test3-1",
					Address: "10.0.0.3:10003",
					Metadata: map[string]string{
						"foo3": "bar3",
					},
				},
			},
		},
	}

	testFn := func(service, s *registry.Service) {
		if s == nil {
			t.Fatalf("Expected one result for %s got nil", service.Name)

		}

		if s.Name != service.Name {
			t.Fatalf("Expected name %s got %s", service.Name, s.Name)
		}

		if s.Version != service.Version {
			t.Fatalf("Expected version %s got %s", service.Version, s.Version)
		}

		if len(s.Nodes) != 1 {
			t.Fatalf("Expected 1 node, got %d", len(s.Nodes))
		}

		node := s.Nodes[0]

		if node.Id != service.Nodes[0].Id {
			t.Fatalf("Expected node id %s got %s", service.Nodes[0].Id, node.Id)
		}

		if node.Address != service.Nodes[0].Address {
			t.Fatalf("Expected node address %s got %s", service.Nodes[0].Address, node.Address)
		}
	}

	travis := os.Getenv("TRAVIS")

	var opts []registry.Option

	if travis == "true" {
		opts = append(opts, registry.Timeout(time.Millisecond*100))
	}

	// new registry
	r := NewRegistry(opts...)

	w, err := r.Watch()
	if err != nil {
		t.Fatal(err)
	}
	defer w.Stop()

	for _, service := range testData {
		// register service
		if err := r.Register(service); err != nil {
			t.Fatal(err)
		}

		for {
			res, err := w.Next()
			if err != nil {
				t.Fatal(err)
			}

			if res.Service.Name != service.Name {
				continue
			}

			if res.Action != "create" {
				t.Fatalf("Expected create event got %s for %s", res.Action, res.Service.Name)
			}

			testFn(service, res.Service)
			break
		}

		// deregister
		if err := r.Deregister(service); err != nil {
			t.Fatal(err)
		}

		for {
			res, err := w.Next()
			if err != nil {
				t.Fatal(err)
			}

			if res.Service.Name != service.Name {
				continue
			}

			if res.Action != "delete" {
				continue
			}

			testFn(service, res.Service)
			break
		}
	}
}
