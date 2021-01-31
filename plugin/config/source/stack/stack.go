package stack

import (
	"context"

	"github.com/stack-labs/stack-rpc/client"
	log "github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/pkg/config/source"
	proto "github.com/stack-labs/stack-rpc/plugin/config/source/stack/proto"
)

var (
	DefaultPath        = "/stack"
	DefaultServiceName = "stack.rpc.config"
)

type mucpSource struct {
	serviceName string
	path        string
	opts        source.Options
	client      proto.SourceService
}

func (m *mucpSource) Read() (set *source.ChangeSet, err error) {
	req, err := m.client.Read(context.Background(), &proto.ReadRequest{Path: m.path})
	if err != nil {
		return nil, err
	}

	return toChangeSet(req.ChangeSet), nil
}

func (m *mucpSource) Watch() (w source.Watcher, err error) {
	stream, err := m.client.Watch(context.Background(), &proto.WatchRequest{Path: m.path})
	if err != nil {
		log.Error("watch err: ", err)
		return
	}
	return newWatcher(stream)
}

// Write is unsupported
func (m *mucpSource) Write(cs *source.ChangeSet) error {
	return nil
}

func (m *mucpSource) String() string {
	return "stack"
}

func NewSource(opts ...source.Option) source.Source {
	var options source.Options
	for _, o := range opts {
		o(&options)
	}

	addr := DefaultServiceName
	path := DefaultPath
	var cli client.Client

	if options.Context != nil {
		a, ok := options.Context.Value(serviceNameKey{}).(string)
		if ok {
			addr = a
		}
		p, ok := options.Context.Value(pathKey{}).(string)
		if ok {
			path = p
		}

		c, ok := options.Context.Value(clientKey{}).(client.Client)
		if ok {
			cli = c
		}
	}

	if cli == nil {
		log.Errorf("create new stack source error, needs a client")
	}

	s := &mucpSource{
		serviceName: addr,
		path:        path,
		opts:        options,
		client:      proto.NewSourceService(addr, cli),
	}

	return s
}
