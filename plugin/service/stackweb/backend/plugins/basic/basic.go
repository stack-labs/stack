package basic

import (
	"sync"

	"github.com/stack-labs/stack/plugin/service/stackweb/plugins"
)

func init() {
	m = &basicModule{
		name: "basic",
		path: "/b",
	}

	plugins.Register(m)
}

var (
	m *basicModule

	// Default address to bind to
	GatewayNamespaces  = []string{"stack.stackway.api"}
	WebNamespacePrefix = []string{"stack.web"}
)

// basicModule includes web, registry, CLI, Stats submodules.
type basicModule struct {
	name string
	path string
	sync.RWMutex
	api *api
}

func (m *basicModule) Name() string {
	return m.name
}

func (m *basicModule) Path() string {
	return m.path
}

func (m *basicModule) Init(opts ...plugins.Option) error {
	options := &plugins.Options{}
	for _, opt := range opts {
		opt(options)
	}

	m.api = &api{
		service: options.Service,
	}
	return nil
}

func (m *basicModule) Handlers() (mp map[string]*plugins.Handler) {
	m.Lock()
	defer m.Unlock()

	mp = make(map[string]*plugins.Handler)
	mp["/services"] = &plugins.Handler{
		Func:   m.api.services,
		Method: []string{"GET"},
	}

	mp["/micro-services"] = &plugins.Handler{
		Func:   m.api.microServices,
		Method: []string{"GET"},
	}

	mp["/service"] = &plugins.Handler{
		Func:   m.api.handler,
		Method: []string{"GET"},
	}

	mp["/api-gateway-services"] = &plugins.Handler{
		Func:   m.api.apiGatewayServices,
		Method: []string{"GET"},
	}

	mp["/service-details"] = &plugins.Handler{
		Func:   m.api.serviceDetails,
		Method: []string{"GET"},
	}

	mp["/stats"] = &plugins.Handler{
		Func:   m.api.stats,
		Method: []string{"GET"},
	}

	mp["/web-services"] = &plugins.Handler{
		Func:   m.api.webServices,
		Method: []string{"GET"},
	}

	mp["/rpc"] = &plugins.Handler{
		Func:   m.api.rpc,
		Method: []string{"POST"},
	}

	mp["/health"] = &plugins.Handler{
		Func:   m.api.health,
		Method: []string{"GET"},
	}

	return
}
