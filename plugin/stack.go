package plugin

import (
	"github.com/stack-labs/stack-rpc/broker"
	"github.com/stack-labs/stack-rpc/broker/http"
	"github.com/stack-labs/stack-rpc/broker/memory"
	"github.com/stack-labs/stack-rpc/broker/nats"
	brokerSrv "github.com/stack-labs/stack-rpc/broker/service"
	"github.com/stack-labs/stack-rpc/client"
	cgrpc "github.com/stack-labs/stack-rpc/client/grpc"
	cmucp "github.com/stack-labs/stack-rpc/client/mucp"
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/client/selector/dns"
	"github.com/stack-labs/stack-rpc/client/selector/router"
	"github.com/stack-labs/stack-rpc/client/selector/static"
	"github.com/stack-labs/stack-rpc/registry"
	"github.com/stack-labs/stack-rpc/registry/etcd"
	"github.com/stack-labs/stack-rpc/registry/mdns"
	rmem "github.com/stack-labs/stack-rpc/registry/memory"
	regSrv "github.com/stack-labs/stack-rpc/registry/service"
	"github.com/stack-labs/stack-rpc/runtime"
	"github.com/stack-labs/stack-rpc/runtime/kubernetes"
	"github.com/stack-labs/stack-rpc/server"
	sgrpc "github.com/stack-labs/stack-rpc/server/grpc"
	smucp "github.com/stack-labs/stack-rpc/server/mucp"
	"github.com/stack-labs/stack-rpc/transport"
	tgrpc "github.com/stack-labs/stack-rpc/transport/grpc"
	thttp "github.com/stack-labs/stack-rpc/transport/http"
	tmem "github.com/stack-labs/stack-rpc/transport/memory"
	"github.com/stack-labs/stack-rpc/transport/quic"
)

var (
	DefaultBrokers = map[string]func(...broker.Option) broker.Broker{
		"stack.rpc.broker": brokerSrv.NewBroker,
		"service":          brokerSrv.NewBroker,
		"http":             http.NewBroker,
		"memory":           memory.NewBroker,
		"nats":             nats.NewBroker,
	}

	DefaultClients = map[string]func(...client.Option) client.Client{
		"rpc":  client.NewClient,
		"mucp": cmucp.NewClient,
		"grpc": cgrpc.NewClient,
	}

	DefaultRegistries = map[string]func(...registry.Option) registry.Registry{
		"stack.rpc.registry": regSrv.NewRegistry,
		"service":            regSrv.NewRegistry,
		"etcd":               etcd.NewRegistry,
		"mdns":               mdns.NewRegistry,
		"memory":             rmem.NewRegistry,
	}

	DefaultSelectors = map[string]func(...selector.Option) selector.Selector{
		"default": selector.NewSelector,
		"dns":     dns.NewSelector,
		"cache":   selector.NewSelector,
		"router":  router.NewSelector,
		"static":  static.NewSelector,
	}

	DefaultServers = map[string]func(...server.Option) server.Server{
		"rpc":  server.NewServer,
		"mucp": smucp.NewServer,
		"grpc": sgrpc.NewServer,
	}

	DefaultTransports = map[string]func(...transport.Option) transport.Transport{
		"memory": tmem.NewTransport,
		"http":   thttp.NewTransport,
		"grpc":   tgrpc.NewTransport,
		"quic":   quic.NewTransport,
	}

	DefaultRuntimes = map[string]func(...runtime.Option) runtime.Runtime{
		"local":      runtime.NewRuntime,
		"kubernetes": kubernetes.NewRuntime,
	}
)
