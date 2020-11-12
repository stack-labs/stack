package stack

type stackConfig struct {
	Broker               string   `json:"broker"`
	BrokerAddress        string   `json:"broker_address"`
	Client               string   `json:"client"`
	ClientPoolSize       int      `json:"client_pool_size"`
	ClientPoolTTL        string   `json:"client_pool_ttl"`
	ClientRequestTimeout string   `json:"client_request_timeout"`
	ClientRetries        int      `json:"client_retries"`
	Profile              string   `json:"profile"`
	RegisterInterval     int      `json:"register_interval"`
	RegisterTTL          int      `json:"register_ttl"`
	Registry             string   `json:"registry"`
	RegistryAddress      string   `json:"registry_address"`
	Runtime              string   `json:"runtime"`
	Selector             string   `json:"selector"`
	Server               string   `json:"server"`
	ServerAddress        string   `json:"server_address"`
	ServerAdvertise      string   `json:"server_advertise"`
	ServerID             string   `json:"server_id"`
	ServerMetadata       []string `json:"server_metadata"`
	ServerName           string   `json:"server_name"`
	ServerVersion        string   `json:"server_version"`
	Transport            string   `json:"transport"`
	TransportAddress     string   `json:"transport_address"`
}

func newDefaultConfig() *stackConfig {
	return &stackConfig{
		ClientRetries:    1,
		RegisterInterval: 30,
		RegisterTTL:      60,
	}
}
