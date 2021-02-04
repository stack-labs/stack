package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/stack-labs/stack"
	"github.com/stack-labs/stack/plugin/service/stackway/handler"
	"github.com/stack-labs/stack/plugin/service/stackway/helper"
	"github.com/stack-labs/stack/plugin/service/stackway/plugin"
	gwServer "github.com/stack-labs/stack/plugin/service/stackway/server"
	ahandler "github.com/stack-labs/stack/api/handler"
	aapi "github.com/stack-labs/stack/api/handler/api"
	"github.com/stack-labs/stack/api/handler/event"
	ahttp "github.com/stack-labs/stack/api/handler/http"
	arpc "github.com/stack-labs/stack/api/handler/rpc"
	"github.com/stack-labs/stack/api/handler/web"
	"github.com/stack-labs/stack/api/resolver"
	"github.com/stack-labs/stack/api/resolver/grpc"
	"github.com/stack-labs/stack/api/resolver/host"
	"github.com/stack-labs/stack/api/resolver/path"
	rrstack "github.com/stack-labs/stack/api/resolver/stack"
	"github.com/stack-labs/stack/api/router"
	regRouter "github.com/stack-labs/stack/api/router/registry"
	apiServer "github.com/stack-labs/stack/api/server"
	"github.com/stack-labs/stack/api/server/acme"
	"github.com/stack-labs/stack/api/server/acme/autocert"
	httpapi "github.com/stack-labs/stack/api/server/http"
	"github.com/stack-labs/stack/service"
	"github.com/stack-labs/stack/util/log"
)

type config struct {
	Server   *server   `json:"server"`
	Stackway *stackway `json:"stackway"`
}

type server struct {
	Address string `json:"address"`
}

type stackway struct {
	Address      string      `json:"address"`
	Handler      string      `json:"handler"`
	Resolver     string      `json:"resolver"`
	RPCPath      string      `json:"rpc_path"`
	APIPath      string      `json:"api_path"`
	ProxyPath    string      `json:"proxy_path"`
	Namespace    string      `json:"namespace"`
	HeaderPrefix string      `json:"header_prefix"`
	EnableRPC    bool        `json:"enable_rpc"`
	EnableACME   bool        `json:"enable_acme"`
	EnableTLS    bool        `json:"enable_tls"`
	ACME         *acmeConfig `json:"acme"`
	TLS          *helper.TLS `json:"tls"`
}

type acmeConfig struct {
	Provider          string   `json:"provider"`
	ChallengeProvider string   `json:"challenge_provider"`
	CA                string   `json:"ca"`
	Hosts             []string `json:"hosts"`
}

func newDefaultConfig() *config {
	return &config{
		Server: &server{
			Address: ":8080",
		},
		Stackway: &stackway{
			Handler:      "meta",
			Resolver:     "stack",
			RPCPath:      "/rpc",
			APIPath:      "/",
			ProxyPath:    "/{service:[a-zA-Z0-9]+}",
			Namespace:    "stack.rpc.api",
			HeaderPrefix: "X-Stack-",
			EnableRPC:    false,
			ACME: &acmeConfig{
				Provider:          "autocert",
				ChallengeProvider: "cloudflare",
				CA:                acme.LetsEncryptProductionCA,
			},
		},
	}
}

type httpServer struct {
	svc service.Service
	api apiServer.Server
}

func (s *httpServer) Options() []service.Option {
	opts := Options()
	opts = append(
		opts,
		stack.Server(
			gwServer.NewServer(gwServer.HookServer(s)),
		),
	)

	return opts
}

func (s *httpServer) Start() error {
	svc := s.svc
	cfg := svc.Options().Config
	conf := newDefaultConfig()
	if cfg != nil {
		if c := cfg.Get("stack"); c != nil {
			if err := c.Scan(conf); err != nil {
				return err
			}
		}
	}

	log.Debugf("stack config: %v", string(cfg.Bytes()))

	gwConf := conf.Stackway
	address := conf.Server.Address
	if len(gwConf.Address) > 0 {
		address = gwConf.Address
	}

	// Init plugins
	for _, p := range plugin.Plugins() {
		_ = p.Init(cfg)
	}

	// Init API
	var opts []apiServer.Option

	if gwConf.EnableACME {
		opts = append(opts, apiServer.EnableACME(true))
		opts = append(opts, apiServer.ACMEHosts(gwConf.ACME.Hosts...))
		switch gwConf.ACME.Provider {
		case "autocert":
			opts = append(opts, apiServer.ACMEProvider(autocert.New()))
		default:
			log.Fatalf("%s is not s valid ACME provider\n", gwConf.ACME.Provider)
		}
	} else if gwConf.EnableTLS {
		config, err := helper.TLSConfig(gwConf.TLS)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		opts = append(opts, apiServer.EnableTLS(true))
		opts = append(opts, apiServer.TLSConfig(config))
	}

	// create the router
	var h http.Handler
	r := mux.NewRouter()
	h = r

	// return version and list of services
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		helper.ServeCORS(w, r)

		if r.Method == "OPTIONS" {
			return
		}

		// TODO index custom
		response := fmt.Sprintf(`{"version": "%s"}`, svc.Server().Options().Version)
		_, _ = w.Write([]byte(response))
	})

	// strip favicon.ico
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})

	// srvOpts = append(srvOpts, stack.Name(Name))
	// if i := time.Duration(ctx.GlobalInt("register_ttl")); i > 0 {
	// 	srvOpts = append(srvOpts, stack.RegisterTTL(i*time.Second))
	// }
	// if i := time.Duration(ctx.GlobalInt("register_interval")); i > 0 {
	// 	srvOpts = append(srvOpts, stack.RegisterInterval(i*time.Second))
	// }

	// initialise svc
	// svc := stack.NewService(srvOpts...)
	// register rpc handler
	if gwConf.EnableRPC {
		log.Logf("Registering RPC Handler at %s", gwConf.RPCPath)
		r.Handle(gwConf.RPCPath, handler.NewRPCHandlerFunc(svc.Options()))
	}

	// resolver options
	ropts := []resolver.Option{
		resolver.WithNamespace(gwConf.Namespace),
		resolver.WithHandler(gwConf.Handler),
	}

	// default resolver
	rr := rrstack.NewResolver(ropts...)

	switch gwConf.Resolver {
	case "host":
		rr = host.NewResolver(ropts...)
	case "path":
		rr = path.NewResolver(ropts...)
	case "grpc":
		rr = grpc.NewResolver(ropts...)
	}

	switch gwConf.Handler {
	case "rpc":
		log.Logf("Registering API RPC Handler at %s", gwConf.APIPath)
		rt := regRouter.NewRouter(
			router.WithNamespace(gwConf.Namespace),
			router.WithHandler(arpc.Handler),
			router.WithResolver(rr),
			router.WithRegistry(svc.Options().Registry),
		)
		rp := arpc.NewHandler(
			ahandler.WithNamespace(gwConf.Namespace),
			ahandler.WithRouter(rt),
			ahandler.WithService(svc),
		)
		r.PathPrefix(gwConf.APIPath).Handler(rp)
	case "api":
		log.Logf("Registering API Request Handler at %s", gwConf.APIPath)
		rt := regRouter.NewRouter(
			router.WithNamespace(gwConf.Namespace),
			router.WithHandler(aapi.Handler),
			router.WithResolver(rr),
			router.WithRegistry(svc.Options().Registry),
		)
		ap := aapi.NewHandler(
			ahandler.WithNamespace(gwConf.Namespace),
			ahandler.WithRouter(rt),
			ahandler.WithService(svc),
		)
		r.PathPrefix(gwConf.APIPath).Handler(ap)
	case "event":
		log.Logf("Registering API Event Handler at %s", gwConf.APIPath)
		rt := regRouter.NewRouter(
			router.WithNamespace(gwConf.Namespace),
			router.WithHandler(event.Handler),
			router.WithResolver(rr),
			router.WithRegistry(svc.Options().Registry),
		)
		ev := event.NewHandler(
			ahandler.WithNamespace(gwConf.Namespace),
			ahandler.WithRouter(rt),
			ahandler.WithService(svc),
		)
		r.PathPrefix(gwConf.APIPath).Handler(ev)
	case "http", "proxy":
		log.Logf("Registering API HTTP Handler at %s", gwConf.ProxyPath)
		rt := regRouter.NewRouter(
			router.WithNamespace(gwConf.Namespace),
			router.WithHandler(ahttp.Handler),
			router.WithResolver(rr),
			router.WithRegistry(svc.Options().Registry),
		)
		ht := ahttp.NewHandler(
			ahandler.WithNamespace(gwConf.Namespace),
			ahandler.WithRouter(rt),
			ahandler.WithService(svc),
		)
		r.PathPrefix(gwConf.ProxyPath).Handler(ht)
	case "web":
		log.Logf("Registering API Web Handler at %s", gwConf.APIPath)
		rt := regRouter.NewRouter(
			router.WithNamespace(gwConf.Namespace),
			router.WithHandler(web.Handler),
			router.WithResolver(rr),
			router.WithRegistry(svc.Options().Registry),
		)
		w := web.NewHandler(
			ahandler.WithNamespace(gwConf.Namespace),
			ahandler.WithRouter(rt),
			ahandler.WithService(svc),
		)
		r.PathPrefix(gwConf.APIPath).Handler(w)
	default:
		log.Logf("Registering API Default Handler at %s", gwConf.APIPath)
		rt := regRouter.NewRouter(
			router.WithNamespace(gwConf.Namespace),
			router.WithResolver(rr),
			router.WithRegistry(svc.Options().Registry),
		)
		r.PathPrefix(gwConf.APIPath).Handler(handler.Meta(svc, rt))
	}

	// reverse wrap handler
	plugins := append(plugin.Plugins(), plugin.Plugins()...)
	for i := len(plugins); i > 0; i-- {
		h = plugins[i-1].Handler()(h)
	}

	// create the server
	api := httpapi.NewServer(address)
	_ = api.Init(opts...)
	api.Handle("/", h)

	s.api = api

	return s.api.Start()
}

func (s *httpServer) Stop() error {
	return s.api.Stop()
}

func NewServer(svc service.Service) *httpServer {
	return &httpServer{svc: svc}
}
