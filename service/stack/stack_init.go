package stack

import (
	"context"
	"fmt"
	br "github.com/stack-labs/stack-rpc/broker"
	cl "github.com/stack-labs/stack-rpc/client"
	sel "github.com/stack-labs/stack-rpc/client/selector"
	"github.com/stack-labs/stack-rpc/plugin"
	ser "github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/service"
	"github.com/stack-labs/stack-rpc/util/log"
	"github.com/stack-labs/stack-rpc/util/wrapper"
)

// Init initialises options. Additionally it calls cmd.Init
// which parses command line flags. cmd.Init is only called
// on first Init.
func (s *stackService) Init(opts ...service.Option) error {
	// process options
	for _, o := range opts {
		o(&s.opts)
	}

	if s.opts.Context == nil {
		s.opts.Context = context.Background()
	}

	if len(s.opts.BeforeInit) > 0 {
		for _, f := range s.opts.BeforeInit {
			err := f(&s.opts)
			if err != nil {
				log.Fatalf("init service err: %s", err)
			}
		}
	}

	// begin init
	if err := s.initComponents(); err != nil {
		log.Fatalf("init service's components err: %s", err)
	}

	return nil
}

func (s *stackService) initComponents() error {
	serverOpts := s.opts.ServerOptions.Options()
	clientOpts := s.opts.ClientOptions.Options()
	transOpts := s.opts.TransportOptions.Options()
	selectorOpts := s.opts.SelectorOptions.Options()
	regOpts := s.opts.RegistryOptions.Options()
	brokerOpts := s.opts.BrokerOptions.Options()
	logOpts := s.opts.LoggerOptions.Options()

	// set Logger
	// only change if we have the logger and type differs
	if len(logOpts.Name) > 0 && s.opts.Logger.String() != logOpts.Name {
		l, ok := plugin.LoggerPlugins[logOpts.Name]
		if !ok {
			return fmt.Errorf("logger [%s] not found", brokerOpts.Name)
		}

		s.opts.Logger = l.New()
	}

	// Set the client
	if len(clientOpts.Protocol) > 0 {
		// only change if we have the client and type differs
		if cl, ok := plugin.ClientPlugins[clientOpts.Protocol]; ok && s.opts.Client.String() != clientOpts.Protocol {
			s.opts.Client = cl.New()
		}
	}

	// Set the server
	if len(serverOpts.Protocol) > 0 {
		// only change if we have the server and type differs
		if ser, ok := plugin.ServerPlugins[serverOpts.Protocol]; ok && s.opts.Server.String() != serverOpts.Protocol {
			s.opts.Server = ser.New()
		}
	}

	// Set the broker
	if len(brokerOpts.Name) > 0 && s.opts.Broker.String() != brokerOpts.Name {
		b, ok := plugin.BrokerPlugins[brokerOpts.Name]
		if !ok {
			return fmt.Errorf("broker [%s] not found", brokerOpts.Name)
		}

		s.opts.Broker = b.New()
	}

	// Set the registry
	if len(regOpts.Name) > 0 && s.opts.Registry.String() != regOpts.Name {
		r, ok := plugin.RegistryPlugins[regOpts.Name]
		if !ok {
			return fmt.Errorf("registry [%s] not found", regOpts.Name)
		}

		s.opts.Registry = r.New()
	}

	// Set the selector
	if len(selectorOpts.Name) > 0 && s.opts.Selector.String() != selectorOpts.Name {
		sl, ok := plugin.SelectorPlugins[selectorOpts.Name]
		if !ok {
			return fmt.Errorf("selector [%s] not found", selectorOpts.Name)
		}

		s.opts.Selector = sl.New()
	}

	// Set the transport
	if len(transOpts.Name) > 0 && s.opts.Transport.String() != transOpts.Name {
		t, ok := plugin.TransportPlugins[transOpts.Name]
		if !ok {
			return fmt.Errorf("transport [%s] not found", transOpts.Name)
		}

		s.opts.Transport = t.New()
	}

	// set client name
	if len(clientOpts.Name) != 0 {
		s.opts.ClientOptions = append(s.opts.ClientOptions, cl.Name(clientOpts.Name))
	} else {
		clientName := fmt.Sprintf("%s-client", s.Name())
		s.opts.ClientOptions = append(s.opts.ClientOptions, cl.Name(clientName))
	}

	if len(serverOpts.Name) != 0 {
		s.opts.ServerOptions = append(s.opts.ServerOptions, ser.Name(serverOpts.Name))
	} else {
		// serverName := fmt.Sprintf("%s-server", s.Name())
		s.opts.ServerOptions = append(s.opts.ServerOptions, ser.Name(s.Name()))
	}

	s.opts.ServerOptions = append(s.opts.ServerOptions, ser.Transport(s.opts.Transport), ser.Broker(s.opts.Broker), ser.Registry(s.opts.Registry))
	s.opts.ClientOptions = append(s.opts.ClientOptions, cl.Transport(s.opts.Transport), cl.Broker(s.opts.Broker), cl.Registry(s.opts.Registry), cl.Selector(s.opts.Selector))
	s.opts.SelectorOptions = append(s.opts.SelectorOptions, sel.Registry(s.opts.Registry))
	s.opts.BrokerOptions = append(s.opts.BrokerOptions, br.Registry(s.opts.Registry))

	if err := s.opts.Auth.Init(s.opts.AuthOptions...); err != nil {
		return fmt.Errorf("Error configuring auth: %v ", err)
	}

	if err := s.opts.Logger.Init(s.opts.LoggerOptions...); err != nil {
		return fmt.Errorf("Error configuring logger: %s ", err)
	}

	if err := s.opts.Broker.Init(s.opts.BrokerOptions...); err != nil {
		return fmt.Errorf("Error configuring broker: %s ", err)
	}

	if err := s.opts.Registry.Init(s.opts.RegistryOptions...); err != nil {
		return fmt.Errorf("Error configuring registry: %s ", err)
	}

	if err := s.opts.Transport.Init(s.opts.TransportOptions...); err != nil {
		return fmt.Errorf("Error configuring transport: %s ", err)
	}

	if err := s.opts.Transport.Init(s.opts.TransportOptions...); err != nil {
		return fmt.Errorf("Error configuring transport: %s ", err)
	}

	if err := s.opts.Selector.Init(s.opts.SelectorOptions...); err != nil {
		return fmt.Errorf("Error configuring selector: %s ", err)
	}

	// wrap client to inject From-Service header on any calls
	// todo wrap not here
	s.opts.Client = wrapper.FromService(s.Name(), s.opts.Client)
	if err := s.opts.Client.Init(s.opts.ClientOptions...); err != nil {
		return fmt.Errorf("Error configuring client: %v ", err)
	}

	if err := s.opts.Server.Init(s.opts.ServerOptions...); err != nil {
		return fmt.Errorf("Error configuring server: %v ", err)
	}

	return nil
}
