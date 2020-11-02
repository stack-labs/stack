package grpc

import (
	"context"
	"crypto/tls"
	"sync"
	"testing"

	"github.com/stack-labs/stack-rpc"
	"github.com/stack-labs/stack-rpc/registry/memory"
	hello "github.com/stack-labs/stack-rpc/service/grpc/proto"
	mls "github.com/stack-labs/stack-rpc/util/tls"
)

type testHandler struct{}

func (t *testHandler) Call(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	rsp.Msg = "Hello " + req.Name
	return nil
}

func TestGRPCService(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create memory registry
	r := memory.NewRegistry()

	// create GRPC service
	service := NewService(
		stack.Name("test.service"),
		stack.Registry(r),
		stack.AfterStart(func() error {
			wg.Done()
			return nil
		}),
		stack.Context(ctx),
	)

	// register test handler
	hello.RegisterTestHandler(service.Server(), &testHandler{})

	// run service
	go func() {
		if err := service.Run(); err != nil {
			t.Fatal(err)
		}
	}()

	// wait for start
	wg.Wait()

	// create client
	test := hello.NewTestService("test.service", service.Client())

	// call service
	rsp, err := test.Call(context.Background(), &hello.Request{
		Name: "John",
	})
	if err != nil {
		t.Fatal(err)
	}

	// check message
	if rsp.Msg != "Hello John" {
		t.Fatalf("unexpected response %s", rsp.Msg)
	}
}

func TestGRPCFunction(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create service
	fn := NewFunction(
		stack.Name("test.function"),
		stack.Registry(memory.NewRegistry()),
		stack.AfterStart(func() error {
			wg.Done()
			return nil
		}),
		stack.Context(ctx),
	)

	// register test handler
	hello.RegisterTestHandler(fn.Server(), &testHandler{})

	// run service
	go fn.Run()

	// wait for start
	wg.Wait()

	// create client
	test := hello.NewTestService("test.function", fn.Client())

	// call service
	rsp, err := test.Call(context.Background(), &hello.Request{
		Name: "John",
	})
	if err != nil {
		t.Fatal(err)
	}

	// check message
	if rsp.Msg != "Hello John" {
		t.Fatalf("unexpected response %s", rsp.Msg)
	}
}

func TestGRPCTLSService(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create memory registry
	r := memory.NewRegistry()

	// create cert
	cert, err := mls.Certificate("test.service")
	if err != nil {
		t.Fatal(err)
	}
	config := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	// create GRPC service
	service := NewService(
		stack.Name("test.service"),
		stack.Registry(r),
		stack.AfterStart(func() error {
			wg.Done()
			return nil
		}),
		stack.Context(ctx),
		// set TLS config
		WithTLS(config),
	)

	// register test handler
	hello.RegisterTestHandler(service.Server(), &testHandler{})

	// run service
	go func() {
		if err := service.Run(); err != nil {
			t.Fatal(err)
		}
	}()

	// wait for start
	wg.Wait()

	// create client
	test := hello.NewTestService("test.service", service.Client())

	// call service
	rsp, err := test.Call(context.Background(), &hello.Request{
		Name: "John",
	})
	if err != nil {
		t.Fatal(err)
	}

	// check message
	if rsp.Msg != "Hello John" {
		t.Fatalf("unexpected response %s", rsp.Msg)
	}
}
