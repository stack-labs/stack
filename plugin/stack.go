package plugin

var (
	ConfigPlugins    = map[string]ConfigPlugin{}
	ClientPlugins    = map[string]ClientPlugin{}
	ServerPlugins    = map[string]ServerPlugin{}
	BrokerPlugins    = map[string]BrokerPlugin{}
	TransportPlugins = map[string]TransportPlugin{}
	SelectorPlugins  = map[string]SelectorPlugin{}
	RegistryPlugins  = map[string]RegistryPlugin{}
	LoggerPlugins    = map[string]LoggerPlugin{}
)
