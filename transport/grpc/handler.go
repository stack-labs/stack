package grpc

import (
	"runtime/debug"

	"google.golang.org/grpc/peer"
	"github.com/stack-labs/stack-rpc/errors"
	"github.com/stack-labs/stack-rpc/transport"
	pb "github.com/stack-labs/stack-rpc/transport/grpc/proto"
	"github.com/stack-labs/stack-rpc/util/log"
)

// microTransport satisfies the pb.TransportServer inteface
type microTransport struct {
	addr string
	fn   func(transport.Socket)
}

func (m *microTransport) Stream(ts pb.Transport_StreamServer) error {
	var err error

	sock := &grpcTransportSocket{
		stream: ts,
		local:  m.addr,
	}

	p, ok := peer.FromContext(ts.Context())
	if ok {
		sock.remote = p.Addr.String()
	}

	defer func() {
		if r := recover(); r != nil {
			log.Log(r, string(debug.Stack()))
			sock.Close()
			err = errors.InternalServerError("stack.rpc.transport", "panic recovered: %v", r)
		}
	}()

	// execute socket func
	m.fn(sock)

	return err
}
