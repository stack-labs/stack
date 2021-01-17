package plugin

import (
	"github.com/stack-labs/stack-rpc/auth"
	"github.com/stack-labs/stack-rpc/auth/token"
	"github.com/stack-labs/stack-rpc/broker"
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/registry"
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/service"
	"github.com/stack-labs/stack-rpc/transport"
)

type Plugin interface {
	Name() string
}

type ServicePlugin interface {
	Plugin
	Options() []service.Option
	New(...service.Option) service.Service
}

type LoggerPlugin interface {
	Plugin
	Options() []logger.Option
	New(...logger.Option) logger.Logger
}

type BrokerPlugin interface {
	Plugin
	Options() []broker.Option
	New(...broker.Option) broker.Broker
}

type TransportPlugin interface {
	Plugin
	Options() []transport.Option
	New(...transport.Option) transport.Transport
}

type ClientPlugin interface {
	Plugin
	Options() []client.Option
	New(...client.Option) client.Client
}

type ServerPlugin interface {
	Plugin
	Options() []server.Option
	New(...server.Option) server.Server
}

type RegistryPlugin interface {
	Plugin
	Options() []registry.Option
	New(...registry.Option) registry.Registry
}

type SelectorPlugin interface {
	Plugin
	Options() []selector.Option
	New(...selector.Option) selector.Selector
}

type ConfigPlugin interface {
	Plugin
	Options() []config.Option
	New(...config.Option) config.Config
}

type AuthPlugin interface {
	Plugin
	Options() []auth.Option
	New(...auth.Option) auth.Auth
}

type AuthTokenProviderPlugin interface {
	Plugin
	Options() []token.Option
	New(...token.Option) token.Provider
}
