package stack

import (
	"context"
	"testing"
	"time"

	"github.com/stack-labs/stack"
	"github.com/stack-labs/stack/registry"
)

func TestServiceTTL(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := stack.NewService(
		stack.Context(ctx),
		stack.Name("stack.rpc.test.ttl"),
		stack.RegisterTTL(time.Second*3),
		// much bigger than ttl
		stack.RegisterInterval(time.Second*20),
	)

	s.Init(stack.AfterStart(func() error {
		// same as ttl
		time.Sleep(time.Second * 3)

		svcs, err := s.Options().Registry.GetService("stack.rpc.test.ttl")
		if err != nil && err != registry.ErrNotFound {
			t.Fatal(err)
		}

		for _, svc := range svcs {
			if len(svc.Nodes) > 0 {
				t.Fatalf("Service %q still has nodes registered", svc.Nodes[0].Address)
			}
		}

		return nil
	}), stack.AfterStart(func() error {
		cancel()
		return nil
	}))

	s.Run()
}
