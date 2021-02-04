module github.com/stack-labs/stack/plugin/registry/etcd

go 1.14

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
	github.com/stack-labs/stack v1.0.1-rc1 => ../../../
)

require (
	github.com/coreos/bbolt v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/mitchellh/hashstructure v1.1.0
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/stack-labs/stack v1.0.1-rc1
	github.com/tmc/grpc-websocket-proxy v0.0.0-20201229170055-e5319fda7802 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.5 // indirect
	go.uber.org/zap v1.16.0 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)
