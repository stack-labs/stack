
GOPATH:=$(shell go env GOPATH)

.PHONY: proto
proto:
	protoc --proto_path=${GOPATH}/src --stack_out=${GOPATH}/src --go_out=:${GOPATH}/src github.com/stack-labs/stack-rpc-plugins/service/stackway/api/proto/api.proto
	protoc --proto_path=${GOPATH}/src --stack_out=${GOPATH}/src --go_out=:${GOPATH}/src github.com/stack-labs/stack-rpc-plugins/service/stackway/test/proto/test.proto

.PHONY: test
test:
	go test -race -cover -v ./...

.PHONY: run
run:
	go run main.go --config stack_config.yml

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-w' -o ./bin/stackway main.go

.PHONY: docker
docker: build
	docker build . -t $(tag)

