module github.com/stack-labs/stack-rpc/examples

go 1.14

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4

	github.com/stack-labs/stack-rpc v1.0.1-rc1 => ../
	github.com/stack-labs/stack-rpc/plugin/config/source/apollo v1.0.1-rc1 => ../plugin/config/source/apollo
	github.com/stack-labs/stack-rpc/plugin/config/source/stack v1.0.1-rc1 => ../plugin/config/source/stack
	github.com/stack-labs/stack-rpc/plugin/logger/logrus v1.0.1-rc1 => ../plugin/logger/logrus
	github.com/stack-labs/stack-rpc/plugin/registry/etcd v1.0.1-rc1 => ../plugin/registry/etcd
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)

require (
	github.com/golang/protobuf v1.4.3
	github.com/stack-labs/stack-rpc v1.0.1-rc1
	github.com/stack-labs/stack-rpc/plugin/config/source/apollo v1.0.1-rc1
	github.com/stack-labs/stack-rpc/plugin/config/source/stack v1.0.1-rc1
	github.com/stack-labs/stack-rpc/plugin/logger/logrus v1.0.1-rc1
	github.com/stack-labs/stack-rpc/plugin/registry/etcd v1.0.1-rc1
	google.golang.org/protobuf v1.25.0
)
