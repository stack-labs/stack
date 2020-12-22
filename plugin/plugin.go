package plugin

import (
	"github.com/stack-labs/stack-rpc/broker"
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/registry"
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/transport"
)

type Options func()

type Plugin interface {
	Name() string
	Type() string
}

type LoggerPlugin interface {
	Plugin
	Config() []logger.Option
}

type BrokerPlugin interface {
	Plugin
	Config() []broker.Option
}

type TransportPlugin interface {
	Plugin
	Config() []transport.Option
}

type ClientPlugin interface {
	Plugin
	Config() []client.Option
}

type ServerPlugin interface {
	Plugin
	Config() []server.Option
}

type RegistryPlugin interface {
	Plugin
	Config() []registry.Option
}

type SelectorPlugin interface {
	Plugin
	Config() []selector.Option
}
