package server

import (
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/stack-labs/stack-rpc/pkg/metadata"
	"github.com/stack-labs/stack-rpc/registry"
	"github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/server/mock"
	"github.com/stack-labs/stack-rpc/util/addr"
	"github.com/stack-labs/stack-rpc/util/log"
	mnet "github.com/stack-labs/stack-rpc/util/net"
)

type stackwayServer struct {
	sync.RWMutex
	running     bool
	opts        server.Options
	handlers    map[string]server.Handler
	subscribers map[string][]server.Subscriber

	exit chan chan error

	// marks the serve as started
	started bool
	// used for first registration
	registered bool
}

var (
	_ server.Server = &stackwayServer{}
)

func newServer(opts ...server.Option) *stackwayServer {
	options := newOptions(opts...)

	return &stackwayServer{
		opts:        options,
		handlers:    make(map[string]server.Handler),
		subscribers: make(map[string][]server.Subscriber),
		exit:        make(chan chan error),
	}
}

func (s *stackwayServer) Options() server.Options {
	s.Lock()
	defer s.Unlock()

	return s.opts
}

func (s *stackwayServer) Init(opts ...server.Option) error {
	s.Lock()
	defer s.Unlock()

	for _, o := range opts {
		o(&s.opts)
	}
	return nil
}

func (s *stackwayServer) Handle(h server.Handler) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.handlers[h.Name()]; ok {
		return errors.New("Handler " + h.Name() + " already exists")
	}
	s.handlers[h.Name()] = h
	return nil
}

func (s *stackwayServer) NewHandler(h interface{}, opts ...server.HandlerOption) server.Handler {
	var options server.HandlerOptions
	for _, o := range opts {
		o(&options)
	}

	return &mock.MockHandler{
		Id:   uuid.New().String(),
		Hdlr: h,
		Opts: options,
	}
}

func (s *stackwayServer) NewSubscriber(topic string, fn interface{}, opts ...server.SubscriberOption) server.Subscriber {
	var options server.SubscriberOptions
	for _, o := range opts {
		o(&options)
	}

	return &mock.MockSubscriber{
		Id:   topic,
		Sub:  fn,
		Opts: options,
	}
}

func (s *stackwayServer) Register() error {
	var err error
	var advt, host, port string

	// parse address for host, port
	config := s.Options()

	// check the advertise address first
	// if it exists then use it, otherwise
	// use the address
	if len(config.Advertise) > 0 {
		advt = config.Advertise
	} else {
		advt = config.Address
	}

	if cnt := strings.Count(advt, ":"); cnt >= 1 {
		// ipv6 address in format [host]:port or ipv4 host:port
		host, port, err = net.SplitHostPort(advt)
		if err != nil {
			return err
		}
	} else {
		host = advt
	}

	address, err := addr.Extract(host)
	if err != nil {
		return err
	}

	// make copy of metadata
	md := make(metadata.Metadata)
	for k, v := range config.Metadata {
		md[k] = v
	}

	// mq-rpc(eg. nats) doesn't need the port. its address is queue name.
	if port != "" {
		address = mnet.HostPort(address, port)
	}

	// register service
	node := &registry.Node{
		Id:       config.Name + "-" + config.Id,
		Address:  address,
		Metadata: md,
	}

	node.Metadata["transport"] = config.Transport.String()
	node.Metadata["broker"] = config.Broker.String()
	node.Metadata["server"] = s.String()
	node.Metadata["registry"] = config.Registry.String()
	node.Metadata["protocol"] = "mucp"

	s.RLock()

	service := &registry.Service{
		Name:    config.Name,
		Version: config.Version,
		Nodes:   []*registry.Node{node},
	}

	// get registered value
	registered := s.registered

	s.RUnlock()

	if !registered {
		log.Logf("Registry [%s] Registering node: %s", config.Registry.String(), node.Id)
	}

	// create registry options
	rOpts := []registry.RegisterOption{registry.RegisterTTL(config.RegisterTTL)}

	if err := config.Registry.Register(service, rOpts...); err != nil {
		return err
	}

	// already registered? don't need to register subscribers
	if registered {
		return nil
	}

	s.Lock()
	defer s.Unlock()

	s.registered = true
	// set what we're advertising
	s.opts.Advertise = address

	return nil
}

func (s *stackwayServer) Deregister() error {
	var err error
	var advt, host, port string

	config := s.Options()

	// check the advertise address first
	// if it exists then use it, otherwise
	// use the address
	if len(config.Advertise) > 0 {
		advt = config.Advertise
	} else {
		advt = config.Address
	}

	if cnt := strings.Count(advt, ":"); cnt >= 1 {
		// ipv6 address in format [host]:port or ipv4 host:port
		host, port, err = net.SplitHostPort(advt)
		if err != nil {
			return err
		}
	} else {
		host = advt
	}

	address, err := addr.Extract(host)
	if err != nil {
		return err
	}

	// mq-rpc(eg. nats) doesn't need the port. its address is queue name.
	if port != "" {
		address = mnet.HostPort(address, port)
	}

	node := &registry.Node{
		Id:      config.Name + "-" + config.Id,
		Address: address,
	}

	service := &registry.Service{
		Name:    config.Name,
		Version: config.Version,
		Nodes:   []*registry.Node{node},
	}

	log.Logf("Registry [%s] Deregistering node: %s", config.Registry.String(), node.Id)
	if err := config.Registry.Deregister(service); err != nil {
		return err
	}

	s.Lock()

	if !s.registered {
		s.Unlock()
		return nil
	}

	s.registered = false

	s.Unlock()
	return nil
}

func (s *stackwayServer) Subscribe(sub server.Subscriber) error {
	s.Lock()
	defer s.Unlock()

	subs := s.subscribers[sub.Topic()]
	subs = append(subs, sub)
	s.subscribers[sub.Topic()] = subs
	return nil
}

func (s *stackwayServer) Start() error {
	s.RLock()
	if s.started {
		s.RUnlock()
		return nil
	}
	s.RUnlock()

	config := s.Options()

	// connect to the broker
	if err := config.Broker.Connect(); err != nil {
		return err
	}

	bname := config.Broker.String()

	log.Logf("Broker [%s] Connected to %s", bname, config.Broker.Address())

	// use RegisterCheck func before register
	if err := s.opts.RegisterCheck(s.opts.Context); err != nil {
		log.Logf("Server %s-%s register check error: %s", config.Name, config.Id, err)
	} else {
		// announce self to the world
		if err = s.Register(); err != nil {
			log.Logf("Server %s-%s register error: %s", config.Name, config.Id, err)
		}
	}

	exit := make(chan bool)

	stopFn := func() error { return nil }
	if s.opts.Context != nil {
		if v := s.opts.Context.Value(hookServerKey{}); v != nil {
			gwServer := v.(hookServer)
			if err := gwServer.Start(); err != nil {
				return err
			}
			stopFn = gwServer.Stop
		}
	}

	go func() {
		t := new(time.Ticker)

		// only process if it exists
		if s.opts.RegisterInterval > time.Duration(0) {
			// new ticker
			t = time.NewTicker(s.opts.RegisterInterval)
		}

		// return error chan
		var ch chan error

	Loop:
		for {
			select {
			// register self on interval
			case <-t.C:
				s.RLock()
				registered := s.registered
				s.RUnlock()
				if err := s.opts.RegisterCheck(s.opts.Context); err != nil && registered {
					log.Logf("Server %s-%s register check error: %s, deregister it", config.Name, config.Id, err)
					// deregister self in case of error
					if err := s.Deregister(); err != nil {
						log.Logf("Server %s-%s deregister error: %s", config.Name, config.Id, err)
					}
				} else {
					if err := s.Register(); err != nil {
						log.Logf("Server %s-%s register error: %s", config.Name, config.Id, err)
					}
				}
			// wait for exit
			case ch = <-s.exit:
				t.Stop()
				close(exit)
				break Loop
			}
		}

		// deregister self
		if err := s.Deregister(); err != nil {
			log.Logf("Server %s-%s deregister error: %s", config.Name, config.Id, err)
		}

		// TODO graceful exit need supported by gateway's api server
		// stop stackway server
		ch <- stopFn()

		// disconnect the broker
		_ = config.Broker.Disconnect()
	}()

	// mark the server as started
	s.Lock()
	s.started = true
	s.Unlock()

	return nil
}

func (s *stackwayServer) Stop() error {
	s.RLock()
	if !s.started {
		s.RUnlock()
		return nil
	}
	s.RUnlock()

	ch := make(chan error)
	s.exit <- ch

	err := <-ch
	s.Lock()
	s.started = false
	s.Unlock()

	return err
}

func (s *stackwayServer) String() string {
	return "stackway"
}

func NewServer(opts ...server.Option) server.Server {
	return newServer(opts...)
}
