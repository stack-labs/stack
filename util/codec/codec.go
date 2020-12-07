package codec

import (
	"github.com/stack-labs/stack-rpc/codec"
	raw "github.com/stack-labs/stack-rpc/codec/bytes"
	"github.com/stack-labs/stack-rpc/codec/grpc"
	"github.com/stack-labs/stack-rpc/codec/json"
	"github.com/stack-labs/stack-rpc/codec/jsonrpc"
	"github.com/stack-labs/stack-rpc/codec/proto"
	"github.com/stack-labs/stack-rpc/codec/protorpc"
)

var (
	DefaultContentType = "application/protobuf"

	DefaultCodecs = map[string]codec.NewCodec{
		"application/grpc":         grpc.NewCodec,
		"application/grpc+json":    grpc.NewCodec,
		"application/grpc+proto":   grpc.NewCodec,
		"application/json":         json.NewCodec,
		"application/json-rpc":     jsonrpc.NewCodec,
		"application/protobuf":     proto.NewCodec,
		"application/proto-rpc":    protorpc.NewCodec,
		"application/octet-stream": raw.NewCodec,
	}

	// TODO: remove legacy codec list
	defaultCodecs = map[string]codec.NewCodec{
		"application/json":         jsonrpc.NewCodec,
		"application/json-rpc":     jsonrpc.NewCodec,
		"application/protobuf":     protorpc.NewCodec,
		"application/proto-rpc":    protorpc.NewCodec,
		"application/octet-stream": protorpc.NewCodec,
	}
)
