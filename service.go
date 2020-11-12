package stack

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/stack-labs/stack-rpc/broker"
	"github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/config"
	"github.com/stack-labs/stack-rpc/debug/profile"
	"github.com/stack-labs/stack-rpc/debug/profile/pprof"
	"github.com/stack-labs/stack-rpc/debug/service/handler"
	"github.com/stack-labs/stack-rpc/plugin"
	"github.com/stack-labs/stack-rpc/registry"
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/transport"
	"github.com/stack-labs/stack-rpc/util/log"
	"github.com/stack-labs/stack-rpc/util/wrapper"
)

type service struct {
	opts Options

	once sync.Once
}

func newService(opts ...Option) Service {
	options := newOptions(opts...)

	// service name
	serviceName := options.Server.Options().Name

	// wrap client to inject From-Service header on any calls
	options.Client = wrapper.FromService(serviceName, options.Client)

	return &service{
		opts: options,
	}
}

func (s *service) Name() string {
	return s.opts.Server.Options().Name
}

// Init initialises options. Additionally it calls cmd.Init
// which parses command line flags. cmd.Init is only called
// on first Init.
func (s *service) Init(opts ...Option) {
	// process options
	for _, o := range opts {
		o(&s.opts)
	}

	s.once.Do(func() {
		// setup the plugins
		for _, p := range strings.Split(os.Getenv("STACK_PLUGIN"), ",") {
			if len(p) == 0 {
				continue
			}

			// load the plugin
			c, err := plugin.Load(p)
			if err != nil {
				log.Fatal(err)
			}

			// initialise the plugin
			if err := plugin.Init(c); err != nil {
				log.Fatal(err)
			}
		}
	})
}

func (s *service) Options() Options {
	return s.opts
}

func (s *service) Client() client.Client {
	return s.opts.Client
}

func (s *service) Server() server.Server {
	return s.opts.Server
}

func (s *service) String() string {
	return "stack"
}

func (s *service) Start() error {
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	if err := s.opts.Server.Start(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStart {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) Stop() error {
	var gerr error

	for _, fn := range s.opts.BeforeStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	if err := s.opts.Server.Stop(); err != nil {
		return err
	}

	if err := s.opts.Config.Close(); err != nil {
		return err
	}

	for _, fn := range s.opts.AfterStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	return gerr
}

func (s *service) Run() error {
	if err := s.opts.Cmd.Init(); err != nil {
		return err
	}

	// init the stack config
	var (
		err      error
		filePath string
	)
	if s.opts.ConfigFile {
		filePath = s.opts.Cmd.ConfigFile()
	}
	if s.opts.Config, err = config.New(filePath, s.opts.Cmd.App(), s.opts.ConfigSource...); err != nil {
		return err
	}

	// load dynamic config
	if err := s.load(); err != nil {
		return err
	}

	// register the debug handler
	if err := s.opts.Server.Handle(
		s.opts.Server.NewHandler(
			handler.DefaultHandler,
			server.InternalHandler(true),
		),
	); err != nil {
		return err
	}

	// start the profiler
	// TODO: set as an option to the service, don't just use pprof
	if prof := os.Getenv("STACK_DEBUG_PROFILE"); len(prof) > 0 {
		service := s.opts.Server.Options().Name
		version := s.opts.Server.Options().Version
		id := s.opts.Server.Options().Id
		profiler := pprof.NewProfile(
			profile.Name(service + "." + version + "." + id),
		)
		if err := profiler.Start(); err != nil {
			return err
		}
		defer profiler.Stop()
	}

	if err := s.Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	if s.opts.Signal {
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	}

	select {
	// wait on kill signal
	case <-ch:
	// wait on context cancel
	case <-s.opts.Context.Done():
	}

	return s.Stop()
}

func (s *service) load() error {
	conf := newDefaultConfig()
	if err := s.opts.Config.Scan(conf); err != nil {
		return err
	}
	log.Infof("stack config: %#v", conf)

	// If flags are set then use them otherwise do nothing
	var serverOpts []server.Option
	var clientOpts []client.Option

	// Set the client
	if len(conf.Client) > 0 {
		// only change if we have the client and type differs
		if cl, ok := plugin.DefaultClients[conf.Client]; ok && s.opts.Client.String() != conf.Client {
			s.opts.Client = cl()
		}
	}

	// Set the server
	if len(conf.Server) > 0 {
		// only change if we have the server and type differs
		if ser, ok := plugin.DefaultServers[conf.Server]; ok && s.opts.Server.String() != conf.Server {
			s.opts.Server = ser()
		}
	}

	// Set the broker
	if len(conf.Broker) > 0 && s.opts.Broker.String() != conf.Broker {
		b, ok := plugin.DefaultBrokers[conf.Broker]
		if !ok {
			return fmt.Errorf("broker %s not found", conf.Broker)
		}

		s.opts.Broker = b()
		serverOpts = append(serverOpts, server.Broker(s.opts.Broker))
		clientOpts = append(clientOpts, client.Broker(s.opts.Broker))
	}

	// Set the registry
	if len(conf.Registry) > 0 && s.opts.Registry.String() != conf.Registry {
		r, ok := plugin.DefaultRegistries[conf.Registry]
		if !ok {
			return fmt.Errorf("registry %s not found", conf.Registry)
		}

		s.opts.Registry = r()
		serverOpts = append(serverOpts, server.Registry(s.opts.Registry))
		clientOpts = append(clientOpts, client.Registry(s.opts.Registry))

		if err := s.opts.Selector.Init(selector.Registry(s.opts.Registry)); err != nil {
			log.Fatalf("Error configuring registry: %v", err)
		}

		clientOpts = append(clientOpts, client.Selector(s.opts.Selector))

		if err := s.opts.Broker.Init(broker.Registry(s.opts.Registry)); err != nil {
			log.Fatalf("Error configuring broker: %v", err)
		}
	}

	// Set the selector
	if len(conf.Selector) > 0 && s.opts.Selector.String() != conf.Selector {
		sel, ok := plugin.DefaultSelectors[conf.Selector]
		if !ok {
			return fmt.Errorf("selector %s not found", conf.Selector)
		}

		s.opts.Selector = sel(selector.Registry(s.opts.Registry))

		// No server option here. Should there be?
		clientOpts = append(clientOpts, client.Selector(s.opts.Selector))
	}

	// Set the transport
	if len(conf.Transport) > 0 && s.opts.Transport.String() != conf.Transport {
		t, ok := plugin.DefaultTransports[conf.Transport]
		if !ok {
			return fmt.Errorf("transport %s not found", conf.Transport)
		}

		s.opts.Transport = t()
		serverOpts = append(serverOpts, server.Transport(s.opts.Transport))
		clientOpts = append(clientOpts, client.Transport(s.opts.Transport))
	}

	// Parse the server options
	metadata := make(map[string]string)
	for _, d := range conf.ServerMetadata {
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

	if len(conf.BrokerAddress) > 0 {
		if err := s.opts.Broker.Init(broker.Addrs(strings.Split(conf.BrokerAddress, ",")...)); err != nil {
			log.Fatalf("Error configuring broker: %v", err)
		}
	}

	if len(conf.RegistryAddress) > 0 {
		if err := s.opts.Registry.Init(registry.Addrs(strings.Split(conf.RegistryAddress, ",")...)); err != nil {
			log.Fatalf("Error configuring registry: %v", err)
		}
	}

	if len(conf.TransportAddress) > 0 {
		if err := s.opts.Transport.Init(transport.Addrs(strings.Split(conf.TransportAddress, ",")...)); err != nil {
			log.Fatalf("Error configuring transport: %v", err)
		}
	}

	if len(conf.ServerName) > 0 {
		serverOpts = append(serverOpts, server.Name(conf.ServerName))
	}

	if len(conf.ServerVersion) > 0 {
		serverOpts = append(serverOpts, server.Version(conf.ServerVersion))
	}

	if len(conf.ServerID) > 0 {
		serverOpts = append(serverOpts, server.Id(conf.ServerID))
	}

	if len(conf.ServerAddress) > 0 {
		serverOpts = append(serverOpts, server.Address(conf.ServerAddress))
	}

	if len(conf.ServerAdvertise) > 0 {
		serverOpts = append(serverOpts, server.Advertise(conf.ServerAdvertise))
	}

	if ttl := time.Duration(conf.RegisterTTL); ttl >= 0 {
		serverOpts = append(serverOpts, server.RegisterTTL(ttl*time.Second))
	}

	if val := time.Duration(conf.RegisterInterval); val >= 0 {
		serverOpts = append(serverOpts, server.RegisterInterval(val*time.Second))
	}

	// client opts
	if conf.ClientRetries >= 0 {
		clientOpts = append(clientOpts, client.Retries(conf.ClientRetries))
	}

	if len(conf.ClientRequestTimeout) > 0 {
		d, err := time.ParseDuration(conf.ClientRequestTimeout)
		if err != nil {
			return fmt.Errorf("failed to parse client_request_timeout: %v", conf.ClientRequestTimeout)
		}
		clientOpts = append(clientOpts, client.RequestTimeout(d))
	}

	if conf.ClientPoolSize > 0 {
		clientOpts = append(clientOpts, client.PoolSize(conf.ClientPoolSize))
	}

	if len(conf.ClientPoolTTL) > 0 {
		d, err := time.ParseDuration(conf.ClientPoolTTL)
		if err != nil {
			return fmt.Errorf("failed to parse client_pool_ttl: %v", conf.ClientPoolTTL)
		}
		clientOpts = append(clientOpts, client.PoolTTL(d))
	}

	// We have some command line opts for the server.
	// Lets set it up
	if len(serverOpts) > 0 {
		if err := s.opts.Server.Init(serverOpts...); err != nil {
			log.Fatalf("Error configuring server: %v", err)
		}
	}

	// Use an init option?
	if len(clientOpts) > 0 {
		if err := s.opts.Client.Init(clientOpts...); err != nil {
			log.Fatalf("Error configuring client: %v", err)
		}
	}

	return nil
}
