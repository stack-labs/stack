package codec

import (
	"github.com/stack-labs/stack/codec"
	raw "github.com/stack-labs/stack/codec/bytes"
	"github.com/stack-labs/stack/codec/grpc"
	"github.com/stack-labs/stack/codec/json"
	"github.com/stack-labs/stack/codec/jsonrpc"
	"github.com/stack-labs/stack/codec/proto"
	"github.com/stack-labs/stack/codec/protorpc"
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
)
