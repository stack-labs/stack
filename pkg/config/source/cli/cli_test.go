package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stack-labs/stack-rpc/pkg/cli"
	"github.com/stack-labs/stack-rpc/pkg/config/source"
)

func test(t *testing.T, withContext bool) {
	var src source.Source

	// setup app
	app := newCmd().App()
	app.Name = "testapp"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "db-host"},
	}

	// with context
	if withContext {
		// set action
		app.Action = func(c *cli.Context) {
			src = WithContext(c)
		}

		// run app
		if err := app.Run([]string{"run", "-db-host", "localhost"}); err != nil {
			t.Error(err)
		}
		// no context
	} else {
		// set args
		os.Args = []string{"run", "-db-host", "localhost"}
		src = NewSource(app)
	}

	// test config
	c, err := src.Read()
	if err != nil {
		t.Error(err)
	}
	if len(c.Data) == 0 {
		t.Fatal()
	}

	t.Log(string(c.Data))

	var actual map[string]interface{}
	if err := json.Unmarshal(c.Data, &actual); err != nil {
		t.Error(err)
	}

	if actual["db-host"] != "localhost" {
		t.Errorf("expected localhost, got %v", actual["name"])
	}
}

func TestCliSource(t *testing.T) {
	// without context
	test(t, false)
}

func TestCliSourceWithContext(t *testing.T) {
	// with context
	test(t, true)
}

func TestCliSource_cmd(t *testing.T) {
	// setup app
	app := newCmd().App()
	app.Name = "testcmd"
	app.Flags = DefaultFlags

	// set args
	os.Args = []string{"run"}
	for _, v := range DefaultFlags {
		os.Args = append(os.Args, fmt.Sprintf("--%s", v.GetName()), "1")
	}
	src := NewSource(app)

	// test config
	c, err := src.Read()
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Data) == 0 {
		t.Fatal()
	}

	t.Log(string(c.Data))

	var conf config
	if err := json.Unmarshal(c.Data, &conf); err != nil {
		t.Error(err)
	}

	// test default
	if conf.Stack.Server.Name != "1" {
		t.Fatal(fmt.Errorf("server name should be [1], not: [%s]", conf.Stack.Server.Name))
	}

	// test the config from cmd
	if conf.Stack.Broker.Address != "1" {
		t.Fatal(fmt.Errorf("broker address should be [1] which is cmd value, not: [%s]", conf.Stack.Broker.Address))
	}

	// test config deep path
	if conf.Stack.Client.Pool.TTL != "1" {
		t.Fatal(fmt.Errorf("client pool's ttl should be [1], not: [%s]", conf.Stack.Client.Pool.TTL))
	}

	// test config root path
	if conf.Stack.Profile != "1" {
		t.Fatal(fmt.Errorf("stack profile should be [1], not: [%s]", conf.Stack.Profile))
	}
}

// region test structs copied from cmd module
type Broker struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type Pool struct {
	Size json.Number `json:"size"`
	TTL  json.Number `json:"ttl"`
}

type ClientRequest struct {
	Retries json.Number `json:"retries"`
	Timeout json.Number `json:"timeout"`
}

type Client struct {
	Protocol string        `json:"protocol"`
	Pool     Pool          `json:"pool"`
	Request  ClientRequest `json:"request"`
}

type Registry struct {
	Address  string      `json:"address"`
	Interval json.Number `json:"interval"`
	Name     string      `json:"name"`
	TTL      json.Number `json:"ttl"`
}

type Metadata map[string]string

type Server struct {
	Address   string   `json:"address"`
	Advertise string   `json:"advertise"`
	ID        string   `json:"id"`
	Metadata  Metadata `json:"metadata"`
	Name      string   `json:"name"`
	Protocol  string   `json:"protocol"`
	Version   string   `json:"version"`
}

type Selector struct {
	Name string `json:"name"`
}

type Transport struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Stack struct {
	Broker    Broker    `json:"broker"`
	Client    Client    `json:"client"`
	Profile   string    `json:"profile"`
	Registry  Registry  `json:"registry"`
	Runtime   string    `json:"runtime"`
	Server    Server    `json:"server"`
	Selector  Selector  `json:"selector"`
	Transport Transport `json:"transport"`
}

type config struct {
	Stack Stack `json:"stack"`
}

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
			// todo good name
			Alias: "stack_client_protocol",
		},
		cli.StringFlag{
			Name:   "client_request_timeout",
			EnvVar: "STACK_CLIENT_REQUEST_TIMEOUT",
			Usage:  "Sets the client request timeout. e.g 500ms, 5s, 1m. Default: 5s",
			Alias:  "stack_client_request_timeout",
		},
		cli.IntFlag{
			Name:   "client_retries",
			EnvVar: "STACK_CLIENT_RETRIES",
			Value:  1,
			Usage:  "Sets the client retries. Default: 1",
			Alias:  "stack_client_retries",
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
			// todo good name
			Alias: "stack_server_protocol",
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
			// todo slice Alias
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
			Name:   "config",
			EnvVar: "STACK_CONFIG",
			Usage:  "config file",
			Value:  "/opt/config.yml",
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

type Option func(o *Options)

type Options struct {
	// For the Command Line itself
	Name        string
	Description string
	Version     string
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// endregion
