module github.com/stack-labs/stack-rpc/examples

go 1.14

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4

	github.com/stack-labs/stack-rpc v1.0.1-rc1 => ../
	github.com/stack-labs/stack-rpc/plugin/registry/etcd v1.0.1-rc1 => ../plugin/registry/etcd
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)

require (
	github.com/coreos/bbolt v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/golang/protobuf v1.4.3
	github.com/stack-labs/stack-rpc v1.0.1-rc1
	github.com/stack-labs/stack-rpc/plugin/registry/etcd v1.0.1-rc1
	github.com/tmc/grpc-websocket-proxy v0.0.0-20201229170055-e5319fda7802 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	google.golang.org/protobuf v1.25.0
)
