package mucp

import (
	"github.com/stack-labs/stack-rpc/broker/http"
	"github.com/stack-labs/stack-rpc/codec"
	"github.com/stack-labs/stack-rpc/registry/mdns"
	"github.com/stack-labs/stack-rpc/server"
	httpt "github.com/stack-labs/stack-rpc/transport/http"
)

func newOptions(opt ...server.Option) server.Options {
	opts := server.Options{
		Codecs:           make(map[string]codec.NewCodec),
		Metadata:         map[string]string{},
		RegisterInterval: server.DefaultRegisterInterval,
		RegisterTTL:      server.DefaultRegisterTTL,
	}

	for _, o := range opt {
		o(&opts)
	}

	if opts.Broker == nil {
		opts.Broker = http.NewBroker()
	}

	if opts.Registry == nil {
		opts.Registry = mdns.NewRegistry()
	}

	if opts.Transport == nil {
		opts.Transport = httpt.NewTransport()
	}

	if opts.RegisterCheck == nil {
		opts.RegisterCheck = server.DefaultRegisterCheck
	}

	if len(opts.Address) == 0 {
		opts.Address = server.DefaultAddress
	}

	if len(opts.Name) == 0 {
		opts.Name = server.DefaultName
	}

	if len(opts.Id) == 0 {
		opts.Id = server.DefaultId
	}

	if len(opts.Version) == 0 {
		opts.Version = server.DefaultVersion
	}

	return opts
}
