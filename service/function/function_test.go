package function

import (
	"context"
	"sync"
	"testing"

	proto "github.com/stack-labs/stack-rpc/debug/service/proto"
	"github.com/stack-labs/stack-rpc/registry/memory"
	"github.com/stack-labs/stack-rpc/service"
	"github.com/stack-labs/stack-rpc/util/test"
)

func TestFunction(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	r := memory.NewRegistry(memory.Services(test.Data))

	// create service
	fn := NewFunction(
		service.Registry(r),
		service.Name("test.function"),
		service.AfterStart(func() error {
			wg.Done()
			return nil
		}),
	)

	// we can't test fn.Init as it parses the command line
	// fn.Init()

	ch := make(chan error, 2)

	go func() {
		// run service
		ch <- fn.Run()
	}()

	// wait for start
	wg.Wait()

	fn.Init()

	// test call debug
	req := fn.Client().NewRequest(
		"test.function",
		"Debug.Health",
		new(proto.HealthRequest),
	)

	rsp := new(proto.HealthResponse)

	err := fn.Client().Call(context.TODO(), req, rsp)
	if err != nil {
		t.Fatal(err)
	}

	if rsp.Status != "ok" {
		t.Fatalf("function response: %s", rsp.Status)
	}

	if err := <-ch; err != nil {
		t.Fatal(err)
	}
}
