package config

import (
	"github.com/stack-labs/stack-rpc/pkg/cli"
	"github.com/stack-labs/stack-rpc/pkg/config"
	"github.com/stack-labs/stack-rpc/pkg/config/reader"
	"github.com/stack-labs/stack-rpc/pkg/config/source"
	cliSource "github.com/stack-labs/stack-rpc/pkg/config/source/cli"
	"github.com/stack-labs/stack-rpc/pkg/config/source/file"
	"github.com/stack-labs/stack-rpc/util/log"
)

type Config interface {
	reader.Values
	Close() error
}

type stackConfig struct {
	config config.Config
}

func New(filePath string, app *cli.App, s ...source.Source) (Config, error) {
	var sources []source.Source
	// need read from config file
	if len(filePath) > 0 {
		log.Info("config read from file:", filePath)
		sources = append(sources, file.NewSource(file.WithPath(filePath)))
	}
	sources = append(sources, cliSource.NewSource(app))
	sources = append(sources, s...)

	c, err := config.NewConfig(config.Storage(true), config.Watch(false))
	if err != nil {
		return nil, err
	}
	if err := c.Load(sources...); err != nil {
		return nil, err
	}

	return &stackConfig{config: c}, nil
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
