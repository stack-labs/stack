package plugin

import (
	"github.com/stack-labs/stack/auth"
	"github.com/stack-labs/stack/auth/token"
	"github.com/stack-labs/stack/broker"
	"github.com/stack-labs/stack/client"
	"github.com/stack-labs/stack/client/selector"
	"github.com/stack-labs/stack/config"
	"github.com/stack-labs/stack/logger"
	"github.com/stack-labs/stack/registry"
	"github.com/stack-labs/stack/server"
	"github.com/stack-labs/stack/service"
	"github.com/stack-labs/stack/transport"
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
