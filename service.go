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
	log "github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/plugin"
	"github.com/stack-labs/stack-rpc/registry"
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/transport"
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
func (s *service) Init(opts ...Option) error {
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

	return nil
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
		log.Errorf("cmd init error: %s", err)
		return err
	}

	// load the all config
	var err error
	if s.opts.Config, err = config.New(
		// todo 合并下列Path与source
		config.FilePath(s.opts.Cmd.ConfigFile()),
		config.Source(s.opts.ConfigSource...),
		config.App(s.opts.Cmd.App()),
	); err != nil {
		log.Errorf("config init error: %s", err)
		return err
	}

	// set config
	config.SetDefaultConfig(s.opts.Config)

	// load service config
	if err := s.loadConfig(); err != nil {
		log.Errorf("load service config error: %s", err)
		return err
	}

	// set Logger
	if err := s.opts.Logger.Init(); err != nil {
		log.Errorf("logger init error: %s", err)
		return err
	}
	log.DefaultLogger = s.opts.Logger

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

func (s *service) loadConfig() error {
	stackConfig := config.GetDefault()
	if err := s.opts.Config.Scan(stackConfig); err != nil {
		return err
	}

	// If flags are set then use them otherwise do nothing
	var serverOpts []server.Option
	var clientOpts []client.Option

	conf := stackConfig.Stack
	// Set the client
	if len(conf.Client.Protocol) > 0 {
		// only change if we have the client and type differs
		if cl, ok := plugin.DefaultClients[conf.Client.Protocol]; ok && s.opts.Client.String() != conf.Client.Protocol {
			s.opts.Client = cl()
		}
	}

	// Set the server
	if len(conf.Server.Protocol) > 0 {
		// only change if we have the server and type differs
		if ser, ok := plugin.DefaultServers[conf.Server.Protocol]; ok && s.opts.Server.String() != conf.Server.Protocol {
			s.opts.Server = ser()
		}
	}

	// Set the broker
	if len(conf.Broker.Name) > 0 && s.opts.Broker.String() != conf.Broker.Name {
		b, ok := plugin.DefaultBrokers[conf.Broker.Name]
		if !ok {
			return fmt.Errorf("broker %s not found", conf.Broker)
		}

		s.opts.Broker = b()
		serverOpts = append(serverOpts, server.Broker(s.opts.Broker))
		clientOpts = append(clientOpts, client.Broker(s.opts.Broker))
	}

	// Set the registry
	if len(conf.Registry.Name) > 0 && s.opts.Registry.String() != conf.Registry.Name {
		r, ok := plugin.DefaultRegistries[conf.Registry.Name]
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
	if len(conf.Selector.Name) > 0 && s.opts.Selector.String() != conf.Selector.Name {
		sel, ok := plugin.DefaultSelectors[conf.Selector.Name]
		if !ok {
			return fmt.Errorf("selector %s not found", conf.Selector)
		}

		s.opts.Selector = sel(selector.Registry(s.opts.Registry))

		// No server option here. Should there be?
		clientOpts = append(clientOpts, client.Selector(s.opts.Selector))
	}

	// Set the transport
	if len(conf.Transport.Name) > 0 && s.opts.Transport.String() != conf.Transport.Name {
		t, ok := plugin.DefaultTransports[conf.Transport.Name]
		if !ok {
			return fmt.Errorf("transport %s not found", conf.Transport)
		}

		s.opts.Transport = t()
		serverOpts = append(serverOpts, server.Transport(s.opts.Transport))
		clientOpts = append(clientOpts, client.Transport(s.opts.Transport))
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

	if len(conf.Broker.Address) > 0 {
		if err := s.opts.Broker.Init(broker.Addrs(strings.Split(conf.Broker.Address, ",")...)); err != nil {
			log.Fatalf("Error configuring broker: %v", err)
		}
	}

	if len(conf.Registry.Address) > 0 {
		if err := s.opts.Registry.Init(registry.Addrs(strings.Split(conf.Registry.Address, ",")...)); err != nil {
			log.Fatalf("Error configuring registry: %v", err)
		}
	}

	if len(conf.Transport.Address) > 0 {
		if err := s.opts.Transport.Init(transport.Addrs(strings.Split(conf.Transport.Address, ",")...)); err != nil {
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
