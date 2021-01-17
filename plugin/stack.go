package plugin

var (
	ServicePlugins           = map[string]ServicePlugin{}
	ConfigPlugins            = map[string]ConfigPlugin{}
	ClientPlugins            = map[string]ClientPlugin{}
	ServerPlugins            = map[string]ServerPlugin{}
	BrokerPlugins            = map[string]BrokerPlugin{}
	TransportPlugins         = map[string]TransportPlugin{}
	SelectorPlugins          = map[string]SelectorPlugin{}
	RegistryPlugins          = map[string]RegistryPlugin{}
	LoggerPlugins            = map[string]LoggerPlugin{}
	AuthPlugins              = map[string]AuthPlugin{}
	AuthTokenProviderPlugins = map[string]AuthTokenProviderPlugin{}
)
