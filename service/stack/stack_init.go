package stack

import (
	"fmt"
	cl "github.com/stack-labs/stack-rpc/client"
	"github.com/stack-labs/stack-rpc/plugin"
	ser "github.com/stack-labs/stack-rpc/server"
	"github.com/stack-labs/stack-rpc/service"
	"github.com/stack-labs/stack-rpc/util/log"
)

// Init initialises options. Additionally it calls cmd.Init
// which parses command line flags. cmd.Init is only called
// on first Init.
func (s *stackService) Init(opts ...service.Option) error {
	// process options
	for _, o := range opts {
		o(&s.opts)
	}

	if len(s.opts.BeforeInit) > 0 {
		for _, f := range s.opts.BeforeInit {
			err := f(&s.opts)
			if err != nil {
				log.Fatalf("init service err: %s", err)
			}
		}
	}

	return nil
}

func (s *stackService) initComponents() error {
	sOpts := conf.Service.Options().opts()

	// serverName := fmt.Sprintf("%s-server", sOpts.Name)
	serverOpts := conf.Server.Options()
	if len(serverOpts.opts().Name) == 0 {
		serverOpts = append(serverOpts, ser.Name(serverName))
	}

	clientName := fmt.Sprintf("%s-client", sOpts.Name)
	clientOpts := conf.Client.Options()
	if len(clientOpts.opts().Name) == 0 {
		clientOpts = append(clientOpts, cl.Name(clientName))
	}

	transOpts := conf.Transport.Options()
	selectorOpts := conf.Selector.Options()
	regOpts := conf.Registry.Options()
	brokerOpts := conf.Broker.Options()
	logOpts := conf.Logger.Options()
	authOpts := conf.Auth.Options()

	// set Logger
	if len(conf.Logger.Name) > 0 {
		// only change if we have the logger and type differs
		if l, ok := plugin.LoggerPlugins[conf.Logger.Name]; ok && (*c.opts.Logger).String() != conf.Logger.Name {
			*c.opts.Logger = l.New()
		}
	}

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
	}

	// Set the registry
	if len(conf.Registry.Name) > 0 && (*c.opts.Registry).String() != conf.Registry.Name {
		r, ok := plugin.RegistryPlugins[conf.Registry.Name]
		if !ok {
			return fmt.Errorf("registry %s not found", conf.Registry.Name)
		}

		*c.opts.Registry = r.New()

		if err := (*c.opts.Selector).Init(sel.Registry(*c.opts.Registry)); err != nil {
			return fmt.Errorf("Error configuring registry: %s ", err)
		}

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

		*c.opts.Selector = sl.New()
	}

	// Set the transport
	if len(conf.Transport.Name) > 0 && (*c.opts.Transport).String() != conf.Transport.Name {
		t, ok := plugin.TransportPlugins[conf.Transport.Name]
		if !ok {
			return fmt.Errorf("transport %s not found", conf.Transport)
		}

		*c.opts.Transport = t.New()
	}

	serverOpts = append(serverOpts, ser.Transport(*c.opts.Transport), ser.Broker(*c.opts.Broker), ser.Registry(*c.opts.Registry))
	clientOpts = append(clientOpts, cl.Transport(*c.opts.Transport), cl.Broker(*c.opts.Broker), cl.Registry(*c.opts.Registry), cl.Selector(*c.opts.Selector))
	selectorOpts = append(selectorOpts, sel.Registry(*c.opts.Registry))

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

	if err = (*c.opts.Transport).Init(transOpts...); err != nil {
		return fmt.Errorf("Error configuring transport: %s ", err)
	}

	if err = (*c.opts.Selector).Init(selectorOpts...); err != nil {
		return fmt.Errorf("Error configuring selector: %s ", err)
	}

	// wrap client to inject From-Service header on any calls
	// todo wrap not here
	*c.opts.Client = wrapper.FromService(serverName, *c.opts.Client)
	if err = (*c.opts.Client).Init(clientOpts...); err != nil {
		return fmt.Errorf("Error configuring client: %v ", err)
	}

	if err = (*c.opts.Auth).Init(authOpts...); err != nil {
		return fmt.Errorf("Error configuring auth: %v ", err)
	}

}
