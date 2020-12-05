// Package cmd is an interface for parsing the command line
package cmd

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/stack-labs/stack-rpc/broker"
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/config"
	log "github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/pkg/cli"
	"github.com/stack-labs/stack-rpc/pkg/config/source"
	cliSource "github.com/stack-labs/stack-rpc/pkg/config/source/cli"
	"github.com/stack-labs/stack-rpc/pkg/config/source/file"
	"github.com/stack-labs/stack-rpc/plugin"
	"github.com/stack-labs/stack-rpc/registry"
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/transport"
)

type Cmd interface {
	// The cli app within this cmd
	App() *cli.App
	// Adds options, parses flags and initialise
	// exits on error
	Init(opts ...Option) error
	// Options set within this command
	Options() Options
	// ConfigFile path. This is not good
	ConfigFile() string
}

type cmd struct {
	opts Options
	app  *cli.App
	conf string
}

var (
	DefaultFlags = []cli.Flag{
		cli.StringFlag{
			Name:   "client",
			EnvVar: "STACK_CLIENT",
			Usage:  "Client for stack-rpc; rpc",
			Alias:  "stack_client_protocol",
		},
		cli.StringFlag{
			Name:   "client_request_timeout",
			EnvVar: "STACK_CLIENT_REQUEST_TIMEOUT",
			Usage:  "Sets the client request timeout. e.g 500ms, 5s, 1m. Default: 5s",
			Alias:  "stack_client_request_timeout",
		},
		cli.IntFlag{
			Name:   "client_request_retries",
			EnvVar: "STACK_CLIENT_REQUEST_RETRIES",
			Value:  1,
			Usage:  "Sets the client retries. Default: 1",
			Alias:  "stack_client_request_retries",
		},
		cli.IntFlag{
			Name:   "client_pool_size",
			EnvVar: "STACK_CLIENT_POOL_SIZE",
			Usage:  "Sets the client connection pool size. Default: 1",
			Alias:  "stack_client_pool_size",
		},
		cli.StringFlag{
			Name:   "client_pool_ttl",
			EnvVar: "STACK_CLIENT_POOL_TTL",
			Usage:  "Sets the client connection pool ttl. e.g 500ms, 5s, 1m. Default: 1m",
			Alias:  "stack_client_pool_ttl",
		},
		cli.IntFlag{
			Name:   "registry_ttl",
			EnvVar: "STACK_REGISTER_TTL",
			Value:  60,
			Usage:  "Register TTL in seconds",
			Alias:  "stack_registry_ttl",
		},
		cli.IntFlag{
			Name:   "registry_interval",
			EnvVar: "STACK_REGISTER_INTERVAL",
			Value:  30,
			Usage:  "Register interval in seconds",
			Alias:  "stack_registry_interval",
		},
		cli.StringFlag{
			Name:   "server",
			EnvVar: "STACK_SERVER",
			Usage:  "Server for stack-rpc; rpc",
			Alias:  "stack_server_protocol",
		},
		cli.StringFlag{
			Name:   "server_name",
			EnvVar: "STACK_SERVER_NAME",
			Usage:  "Name of the server. stack.rpc.srv.example",
			Alias:  "stack_server_name",
		},
		cli.StringFlag{
			Name:   "server_version",
			EnvVar: "STACK_SERVER_VERSION",
			Usage:  "Version of the server. 1.1.0",
			Alias:  "stack_server_version",
		},
		cli.StringFlag{
			Name:   "server_id",
			EnvVar: "STACK_SERVER_ID",
			Usage:  "Id of the server. Auto-generated if not specified",
			Alias:  "stack_server_id",
		},
		cli.StringFlag{
			Name:   "server_address",
			EnvVar: "STACK_SERVER_ADDRESS",
			Usage:  "Bind address for the server. 127.0.0.1:8080",
			Alias:  "stack_server_address",
		},
		cli.StringFlag{
			Name:   "server_advertise",
			EnvVar: "STACK_SERVER_ADVERTISE",
			Usage:  "Used instead of the server_address when registering with discovery. 127.0.0.1:8080",
			Alias:  "stack_server_advertise",
		},
		cli.StringSliceFlag{
			Name:   "server_metadata",
			EnvVar: "STACK_SERVER_METADATA",
			Value:  &cli.StringSlice{},
			Usage:  "A list of key-value pairs defining metadata. version=1.0.0",
			Alias:  "stack_server_metadata",
		},
		cli.StringFlag{
			Name:   "broker",
			EnvVar: "STACK_BROKER",
			Usage:  "Broker for pub/sub. http, nats, rabbitmq",
			Alias:  "stack_broker_name",
		},
		cli.StringFlag{
			Name:   "broker_address",
			EnvVar: "STACK_BROKER_ADDRESS",
			Usage:  "Comma-separated list of broker addresses",
			Alias:  "stack_broker_address",
		},
		cli.StringFlag{
			Name:   "profile",
			Usage:  "Debug profiler for cpu and memory stats",
			EnvVar: "STACK_DEBUG_PROFILE",
			Alias:  "stack_profile",
		},
		cli.StringFlag{
			Name:   "registry",
			EnvVar: "STACK_REGISTRY",
			Usage:  "Registry for discovery. mdns",
			Alias:  "stack_registry_name",
		},
		cli.StringFlag{
			Name:   "registry_address",
			EnvVar: "STACK_REGISTRY_ADDRESS",
			Usage:  "Comma-separated list of registry addresses",
			Alias:  "stack_registry_address",
		},
		cli.StringFlag{
			Name:   "selector",
			EnvVar: "STACK_SELECTOR",
			Usage:  "Selector used to pick nodes for querying",
			Alias:  "stack_selector_name",
		},
		cli.StringFlag{
			Name:   "transport",
			EnvVar: "STACK_TRANSPORT",
			Usage:  "Transport mechanism used; http",
			Alias:  "stack_transport_name",
		},
		cli.StringFlag{
			Name:   "transport_address",
			EnvVar: "STACK_TRANSPORT_ADDRESS",
			Usage:  "Comma-separated list of transport addresses",
			Alias:  "stack_transport_address",
		},
		cli.StringFlag{
			Name:   "logger_level",
			EnvVar: "STACK_LOGGER_LEVEL",
			Usage:  "Logger Level; INFO",
			Alias:  "stack_logger_level",
		},
		cli.StringFlag{
			Name:   "config",
			EnvVar: "STACK_CONFIG",
			Usage:  "config file",
			Alias:  "stack_config",
		},
	}
)

func init() {
	rand.Seed(time.Now().Unix())
	help := cli.HelpPrinter
	cli.HelpPrinter = func(writer io.Writer, templ string, data interface{}) {
		help(writer, templ, data)
		os.Exit(0)
	}
}

func newCmd(opts ...Option) Cmd {
	options := Options{}

	for _, o := range opts {
		o(&options)
	}

	if len(options.Description) == 0 {
		options.Description = "a stack-rpc service"
	}

	cmd := new(cmd)
	cmd.opts = options
	cmd.app = cli.NewApp()
	cmd.app.Name = cmd.opts.Name
	cmd.app.Version = cmd.opts.Version
	cmd.app.Usage = cmd.opts.Description
	cmd.app.Flags = DefaultFlags
	cmd.app.Before = cmd.before
	cmd.app.Action = func(c *cli.Context) {}
	if len(options.Version) == 0 {
		cmd.app.HideVersion = true
	}

	return cmd
}

func (c *cmd) ConfigFile() string {
	return c.conf
}

func (c *cmd) before(ctx *cli.Context) error {
	// set the config file path
	if name := ctx.String("config"); len(name) > 0 {
		c.conf = name
	}

	// need to init config first
	var appendSource []source.Source
	// need read from config file
	if len(c.ConfigFile()) > 0 {
		log.Info("config read from file:", c.ConfigFile())
		configFileSource := file.NewSource(file.WithPath(c.ConfigFile()))
		appendSource = append(appendSource, configFileSource)
	}
	appendSource = append(appendSource, cliSource.NewSource(c.App(), cliSource.Context(c.App().Context())))

	err := (*c.opts.Config).Init(config.Source(appendSource...))
	if err != nil {
		err = fmt.Errorf("init config err: %s", err)
		log.Fatal(err)
		return err
	}

	stackConfig := config.GetDefault()
	if err := (*c.opts.Config).Scan(stackConfig); err != nil {
		return err
	}

	// If flags are set then use them otherwise do nothing
	var serverOpts []server.Option
	var clientOpts []client.Option

	conf := stackConfig.Stack
	// Set the client
	if len(conf.Client.Protocol) > 0 {
		// only change if we have the client and type differs
		if cl, ok := plugin.DefaultClients[conf.Client.Protocol]; ok && (*c.opts.Client).String() != conf.Client.Protocol {
			*c.opts.Client = cl()
		}
	}

	// Set the server
	if len(conf.Server.Protocol) > 0 {
		// only change if we have the server and type differs
		if ser, ok := plugin.DefaultServers[conf.Server.Protocol]; ok && (*c.opts.Server).String() != conf.Server.Protocol {
			*c.opts.Server = ser()
		}
	}

	// Set the broker
	if len(conf.Broker.Name) > 0 && (*c.opts.Broker).String() != conf.Broker.Name {
		b, ok := plugin.DefaultBrokers[conf.Broker.Name]
		if !ok {
			return fmt.Errorf("broker %s not found", conf.Broker)
		}

		*c.opts.Broker = b()
		serverOpts = append(serverOpts, server.Broker(*c.opts.Broker))
		clientOpts = append(clientOpts, client.Broker(*c.opts.Broker))
	}

	// Set the registry
	if len(conf.Registry.Name) > 0 && (*c.opts.Registry).String() != conf.Registry.Name {
		r, ok := plugin.DefaultRegistries[conf.Registry.Name]
		if !ok {
			return fmt.Errorf("registry %s not found", conf.Registry)
		}

		*c.opts.Registry = r()
		serverOpts = append(serverOpts, server.Registry(*c.opts.Registry))
		clientOpts = append(clientOpts, client.Registry(*c.opts.Registry))

		if err := (*c.opts.Selector).Init(selector.Registry(*c.opts.Registry)); err != nil {
			log.Fatalf("Error configuring registry: %v", err)
		}

		clientOpts = append(clientOpts, client.Selector(*c.opts.Selector))

		if err := (*c.opts.Broker).Init(broker.Registry(*c.opts.Registry)); err != nil {
			log.Fatalf("Error configuring broker: %v", err)
		}
	}

	// Set the selector
	if len(conf.Selector.Name) > 0 && (*c.opts.Selector).String() != conf.Selector.Name {
		sel, ok := plugin.DefaultSelectors[conf.Selector.Name]
		if !ok {
			return fmt.Errorf("selector %s not found", conf.Selector)
		}

		*c.opts.Selector = sel(selector.Registry(*c.opts.Registry))

		// No server option here. Should there be?
		clientOpts = append(clientOpts, client.Selector(*c.opts.Selector))
	}

	// Set the transport
	if len(conf.Transport.Name) > 0 && (*c.opts.Transport).String() != conf.Transport.Name {
		t, ok := plugin.DefaultTransports[conf.Transport.Name]
		if !ok {
			return fmt.Errorf("transport %s not found", conf.Transport)
		}

		*c.opts.Transport = t()
		serverOpts = append(serverOpts, server.Transport(*c.opts.Transport))
		clientOpts = append(clientOpts, client.Transport(*c.opts.Transport))
	}

	// Parse the server options
	metadata := make(map[string]string)
	for _, d := range conf.Server.Metadata {
		var key, val string
		parts := strings.Split(d, "=")
		key = parts[0]
		if len(parts) > 1 {
			val = strings.Join(parts[1:], "=")
		}
		metadata[key] = val
	}

	if len(metadata) > 0 {
		serverOpts = append(serverOpts, server.Metadata(metadata))
	}

	// todo we dont need to init so many times
	if len(conf.Broker.Address) > 0 {
		if err := (*c.opts.Broker).Init(broker.Addrs(strings.Split(conf.Broker.Address, ",")...)); err != nil {
			log.Fatalf("Error configuring broker: %v", err)
		}
	} else {
		if err := (*c.opts.Broker).Init(); err != nil {
			log.Fatalf("Error configuring broker: %v", err)
		}
	}

	if len(conf.Registry.Address) > 0 {
		if err := (*c.opts.Registry).Init(registry.Addrs(strings.Split(conf.Registry.Address, ",")...)); err != nil {
			log.Fatalf("Error configuring registry: %v", err)
		}
	} else {
		if err := (*c.opts.Registry).Init(); err != nil {
			log.Fatalf("Error configuring registry: %v", err)
		}
	}

	if len(conf.Transport.Address) > 0 {
		if err := (*c.opts.Transport).Init(transport.Addrs(strings.Split(conf.Transport.Address, ",")...)); err != nil {
			log.Fatalf("Error configuring transport: %v", err)
		}
	} else {
		if err := (*c.opts.Transport).Init(); err != nil {
			log.Fatalf("Error configuring transport: %v", err)
		}
	}

	if len(conf.Server.Name) > 0 {
		serverOpts = append(serverOpts, server.Name(conf.Server.Name))
	}

	if len(conf.Server.Version) > 0 {
		serverOpts = append(serverOpts, server.Version(conf.Server.Version))
	}

	if len(conf.Server.ID) > 0 {
		serverOpts = append(serverOpts, server.Id(conf.Server.ID))
	}

	if len(conf.Server.Address) > 0 {
		serverOpts = append(serverOpts, server.Address(conf.Server.Address))
	}

	if len(conf.Server.Advertise) > 0 {
		serverOpts = append(serverOpts, server.Advertise(conf.Server.Advertise))
	}

	registryTTL, _ := conf.Registry.TTL.Int64()
	if ttl := time.Duration(registryTTL); ttl >= 0 {
		serverOpts = append(serverOpts, server.RegisterTTL(ttl*time.Second))
	}

	registryInterval, _ := conf.Registry.Interval.Int64()
	if val := time.Duration(registryInterval); val >= 0 {
		serverOpts = append(serverOpts, server.RegisterInterval(val*time.Second))
	}

	// client opts
	requestRetries, _ := conf.Client.Request.Retries.Int64()
	if requestRetries >= 0 {
		clientOpts = append(clientOpts, client.Retries(int(requestRetries)))
	}

	if len(conf.Client.Request.Timeout) > 0 {
		d, err := time.ParseDuration(conf.Client.Request.Timeout.String())
		if err != nil {
			return fmt.Errorf("failed to parse client_request_timeout: %v. it shoud be with unit suffix such as 1s, 2m", conf.Client.Request.Timeout.String())
		}
		clientOpts = append(clientOpts, client.RequestTimeout(d))
	}

	if size, _ := conf.Client.Pool.Size.Int64(); size > 0 {
		clientOpts = append(clientOpts, client.PoolSize(int(size)))
	}

	poolTTL := conf.Client.Pool.TTL.String()
	if len(poolTTL) > 0 {
		d, err := time.ParseDuration(poolTTL)
		if err != nil {
			return fmt.Errorf("failed to parse client_pool_ttl: %v. it shoud be with unit suffix such as 1s, 2m", poolTTL)
		}
		clientOpts = append(clientOpts, client.PoolTTL(d))
	}

	// We have some command line opts for the server.
	// Lets set it up
	if len(serverOpts) > 0 {
		if err := (*c.opts.Server).Init(serverOpts...); err != nil {
			log.Fatalf("Error configuring server: %v", err)
		}
	}

	// Use an init option?
	if len(clientOpts) > 0 {
		if err := (*c.opts.Client).Init(clientOpts...); err != nil {
			log.Fatalf("Error configuring client: %v", err)
		}
	}

	return nil
}

func (c *cmd) App() *cli.App {
	return c.app
}

func (c *cmd) Options() Options {
	return c.opts
}

func (c *cmd) Init(opts ...Option) error {
	for _, o := range opts {
		o(&c.opts)
	}
	c.app.Name = c.opts.Name
	c.app.Version = c.opts.Version
	c.app.HideVersion = len(c.opts.Version) == 0
	c.app.Usage = c.opts.Description
	return c.app.Run(os.Args)
}

func NewCmd(opts ...Option) Cmd {
	return newCmd(opts...)
}
