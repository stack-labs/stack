package config

import (
	"github.com/stack-labs/stack-rpc/pkg/config"
	"github.com/stack-labs/stack-rpc/pkg/config/reader"
	"github.com/stack-labs/stack-rpc/pkg/config/source"
)

type Config interface {
	reader.Values

	Config() Value
	Close() error
}

type Value struct {
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

type stackConfig struct {
	config config.Config

	v Value
}

func New(s ...source.Source) (Config, error) {
	c, err := config.NewConfig(config.Storage(true), config.Watch(false))
	if err != nil {
		return nil, err
	}
	if err := c.Load(s...); err != nil {
		return nil, err
	}

	sc := &stackConfig{config: c}
	if err := c.Scan(&sc.v); err != nil {
		return nil, err
	}

	return sc, nil
}

func (c *stackConfig) Config() Value {
	return c.v
}

func (c *stackConfig) Get(path ...string) reader.Value {
	return c.config.Get(path...)
}

func (c *stackConfig) Bytes() []byte {
	return c.config.Bytes()
}

func (c *stackConfig) Map() map[string]interface{} {
	return c.config.Map()
}

func (c *stackConfig) Scan(v interface{}) error {
	return c.config.Scan(v)
}

func (c *stackConfig) Close() error {
	return c.config.Close()
}
