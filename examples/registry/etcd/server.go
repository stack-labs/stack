package main

import (
	"context"
	"time"

	"github.com/stack-labs/stack"
	proto "github.com/stack-labs/stack/examples/proto/service/rpc"
	"github.com/stack-labs/stack/logger"

	_ "github.com/stack-labs/stack/plugin/registry/etcd"
)

// 服务类
type Greeter struct{}

// 实现proto中的Hello接口
func (g Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Greeting = "Hello! " + req.Name
	return nil
}

func main() {
	// 实例化服务，并命名为stack.rpc.registry.etcd
	service := stack.NewService(
		stack.Name("stack.rpc.registry.etcd"),
		stack.RegisterTTL(10*time.Second),
		stack.RegisterInterval(3*time.Second),
	)
	// 初始化服务
	service.Init()

	// 将Greeter注册到服务上
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	// 运行服务
	if err := service.Run(); err != nil {
		logger.Error(err)
	}
}
