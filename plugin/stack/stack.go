package stack

import (
	cfg "github.com/stack-labs/stack/config"
	"github.com/stack-labs/stack/plugin"
)

var options struct {
	Stack struct {
		Service struct {
			GRPC struct {
			} `sc:"grpc"`
			Web struct {
				Enable    bool   `sc:"enable"`
				Address   string `sc:"address"`
				RootPath  string `sc:"root-path"`
				StaticDir string `sc:"static-dir"`
			} `sc:"web"`
		} `sc:"service"`
	} `sc:"stack"`
}

func init() {
	cfg.RegisterOptions(&options)
	plugin.ServerPlugins["grpc"] = &grpcServerPlugin{}
	plugin.ServerPlugins["http"] = &httpServerPlugin{}
	plugin.ServerPlugins["mock"] = &mockServerPlugin{}
	plugin.ServerPlugins["mucp"] = &mucpServerPlugin{}
	plugin.RegistryPlugins["service"] = &serviceRegistryPlugin{}
	plugin.RegistryPlugins["memory"] = &memoryRegistryPlugin{}
	plugin.RegistryPlugins["mdns"] = &mdnsRegistryPlugin{}
	plugin.LoggerPlugins["console"] = &consoleLogPlugin{}
	plugin.ClientPlugins["grpc"] = &grpcClientPlugin{}
	plugin.ClientPlugins["http"] = &httpClientPlugin{}
	plugin.SelectorPlugins["dns"] = &dnsSelectorPlugin{}
	plugin.SelectorPlugins["cache"] = &cacheSelectorPlugin{}
	plugin.SelectorPlugins["static"] = &staticSelectorPlugin{}
	plugin.BrokerPlugins["memory"] = &memoryBrokerPlugin{}
	plugin.BrokerPlugins["http"] = &httpBrokerPlugin{}
	plugin.BrokerPlugins["service"] = &serviceBrokerPlugin{}
	plugin.AuthTokenProviderPlugins["jwt"] = &jwtTokenProviderPlugin{}
	plugin.AuthTokenProviderPlugins["basic"] = &basicTokenProviderPlugin{}
	plugin.AuthPlugins["jwt"] = &jwtAuthPlugin{}
}
