package stack

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
	Stack struct {
		Service struct {
			Name     string `json:"name"`
			RPCPort  int    `json:"rpc-port"`
			HTTPPort int    `json:"http-port"`
		} `json:"service"`
	} `json:"stack"`
}

type stackConfig struct {
	config config.Config

	v Value
}

func newConfig(s ...source.Source) (Config, error) {
	c, err := config.NewConfig(config.Storage(true))
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
