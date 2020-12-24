// Package cmd is an interface for parsing the command line
package cmd

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	br "github.com/stack-labs/stack-rpc/broker"
	cl "github.com/stack-labs/stack-rpc/client"
	sel "github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/config"
	log "github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/pkg/cli"
	"github.com/stack-labs/stack-rpc/pkg/config/source"
	cliSource "github.com/stack-labs/stack-rpc/pkg/config/source/cli"
	"github.com/stack-labs/stack-rpc/pkg/config/source/file"
	"github.com/stack-labs/stack-rpc/plugin"
	ser "github.com/stack-labs/stack-rpc/server"
	"gopkg.in/yaml.v2"
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
			Usage:  "Sets the client connection pool ttl in seconds.",
			Alias:  "stack_client_pool_ttl",
		},
		cli.IntFlag{
			Name:   "server_registry_ttl",
			EnvVar: "STACK_SERVER_REGISTRY_TTL",
			Value:  60,
			Usage:  "Register TTL in seconds",
			Alias:  "stack_server_registry_ttl",
		},
		cli.IntFlag{
			Name:   "server_registry_interval",
			EnvVar: "STACK_SERVER_REGISTRY_INTERVAL",
			Value:  30,
			Usage:  "Register interval in seconds",
			Alias:  "stack_server_registry_interval",
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

	stackStdConfigFile = "stack.yml"
	stackConfig        = StackConfig{}
)

func init() {
	rand.Seed(time.Now().Unix())
	help := cli.HelpPrinter
	cli.HelpPrinter = func(writer io.Writer, templ string, data interface{}) {
		help(writer, templ, data)
		os.Exit(0)
	}

	config.RegisterOptions(&stackConfig)
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

func (c *cmd) before(ctx *cli.Context) (err error) {
	err = c.beforeLoadConfig(ctx)
	if err != nil {
		log.Fatalf("load config in before action err: %s", err)
	}

	err = c.beforeSetupComponents()
	if err != nil {
		log.Fatalf("setup components in before action err: %s", err)
	}

	return nil
}

func (c *cmd) beforeLoadConfig(ctx *cli.Context) (err error) {
	// set the config file path
	if name := ctx.String("config"); len(name) > 0 {
		c.conf = name
	}

	// load the default stack.yml first if exists
	wkDir, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("stack can't access working wkDir: %s", err)
		return
	}

	var appendSource []source.Source
	if len(wkDir) > 0 {
		val := struct {
			Stack struct {
				Includes string `yaml:"includes"`
			} `yaml:"stack"`
		}{}
		stackYmlFile := fmt.Sprintf("%s%s%s", wkDir, string(os.PathSeparator), stackStdConfigFile)
		stdFileSource := file.NewSource(file.WithPath(stackYmlFile))
		appendSource = append(appendSource, stdFileSource)

		set, errN := stdFileSource.Read()
		if errN != nil {
			err = fmt.Errorf("stack read the stack.yml err: %s", errN)
			return
		}
		errN = yaml.Unmarshal(set.Data, &val)
		if errN != nil {
			err = fmt.Errorf("unmarshal stack.yml err: %s", errN)
			return
		}

		if len(val.Stack.Includes) > 0 {
			for _, f := range strings.Split(val.Stack.Includes, ",") {
				log.Infof("load extra config file: %s", f)
				f = strings.TrimSpace(f)
				extraFile := fmt.Sprintf("%s%s%s", wkDir, string(os.PathSeparator), f)
				extraFileSource := file.NewSource(file.WithPath(extraFile))
				appendSource = append(appendSource, extraFileSource)
			}
		}
	}

	// need to init the special config if specified
	if len(c.ConfigFile()) > 0 {
		log.Info("config read from file:", c.ConfigFile())
		configFileSource := file.NewSource(file.WithPath(c.ConfigFile()))
		appendSource = append(appendSource, configFileSource)
	}

	// the last two must be env & cmd line
	appendSource = append(appendSource, cliSource.NewSource(c.App(), cliSource.Context(c.App().Context())))

	err = (*c.opts.Config).Init(config.Source(appendSource...))
	if err != nil {
		err = fmt.Errorf("init config err: %s", err)
		return
	}

	return
}

func (c *cmd) beforeSetupComponents() (err error) {
	conf := stackConfig.Stack

	var serverOpts = conf.Server.Options()
	var clientOpts = conf.Client.Options()
	var transOpts = conf.Transport.Options()
	var regOpts = conf.Registry.Options()
	var brokerOpts = conf.Broker.Options()
	var logOpts = conf.Logger.Options()

	// Set the client
	if len(conf.Client.Protocol) > 0 {
		// only change if we have the client and type differs
		if cl, ok := plugin.ClientPlugins[conf.Client.Protocol]; ok && (*c.opts.Client).String() != conf.Client.Protocol {
			*c.opts.Client = cl.New()
		}
	}

	// Set the server
	if len(conf.Server.Protocol) > 0 {
		// only change if we have the server and type differs
		if ser, ok := plugin.ServerPlugins[conf.Server.Protocol]; ok && (*c.opts.Server).String() != conf.Server.Protocol {
			*c.opts.Server = ser.New()
		}
	}

	// Set the broker
	if len(conf.Broker.Name) > 0 && (*c.opts.Broker).String() != conf.Broker.Name {
		b, ok := plugin.BrokerPlugins[conf.Broker.Name]
		if !ok {
			return fmt.Errorf("broker %s not found", conf.Broker)
		}

		*c.opts.Broker = b.New()
		serverOpts = append(serverOpts, ser.Broker(*c.opts.Broker))
		clientOpts = append(clientOpts, cl.Broker(*c.opts.Broker))
	}

	// Set the registry
	if len(conf.Registry.Name) > 0 && (*c.opts.Registry).String() != conf.Registry.Name {
		r, ok := plugin.RegistryPlugins[conf.Registry.Name]
		if !ok {
			return fmt.Errorf("registry %s not found", conf.Registry.Name)
		}

		*c.opts.Registry = r.New()
		serverOpts = append(serverOpts, ser.Registry(*c.opts.Registry))
		clientOpts = append(clientOpts, cl.Registry(*c.opts.Registry))

		if err := (*c.opts.Selector).Init(sel.Registry(*c.opts.Registry)); err != nil {
			return fmt.Errorf("Error configuring registry: %s ", err)
		}

		clientOpts = append(clientOpts, cl.Selector(*c.opts.Selector))

		if err := (*c.opts.Broker).Init(br.Registry(*c.opts.Registry)); err != nil {
			return fmt.Errorf("Error configuring broker: %s ", err)
		}
	}

	// Set the selector
	if len(conf.Selector.Name) > 0 && (*c.opts.Selector).String() != conf.Selector.Name {
		sl, ok := plugin.SelectorPlugins[conf.Selector.Name]
		if !ok {
			return fmt.Errorf("selector %s not found", conf.Selector)
		}

		*c.opts.Selector = sl.New(sel.Registry(*c.opts.Registry))

		// No server option here. Should there be?
		clientOpts = append(clientOpts, cl.Selector(*c.opts.Selector))
	}

	// Set the transport
	if len(conf.Transport.Name) > 0 && (*c.opts.Transport).String() != conf.Transport.Name {
		t, ok := plugin.TransportPlugins[conf.Transport.Name]
		if !ok {
			return fmt.Errorf("transport %s not found", conf.Transport)
		}

		*c.opts.Transport = t.New()
		serverOpts = append(serverOpts, ser.Transport(*c.opts.Transport))
		clientOpts = append(clientOpts, cl.Transport(*c.opts.Transport))
	}

	if err = (*c.opts.Logger).Init(logOpts...); err != nil {
		return fmt.Errorf("Error configuring logger: %s ", err)
	}

	if err = (*c.opts.Broker).Init(brokerOpts...); err != nil {
		return fmt.Errorf("Error configuring broker: %s ", err)
	}

	if err = (*c.opts.Registry).Init(regOpts...); err != nil {
		return fmt.Errorf("Error configuring registry: %s ", err)
	}

	if err = (*c.opts.Transport).Init(transOpts...); err != nil {
		return fmt.Errorf("Error configuring transport: %s ", err)
	}

	if err = (*c.opts.Server).Init(serverOpts...); err != nil {
		return fmt.Errorf("Error configuring server: %s ", err)
	}

	if err = (*c.opts.Client).Init(clientOpts...); err != nil {
		return fmt.Errorf("Error configuring client: %v ", err)
	}

	return
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
