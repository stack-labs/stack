package config

import (
	"encoding/json"
	"github.com/stack-labs/stack-rpc/pkg/config/source/memory"
	"strings"

	"github.com/stack-labs/stack-rpc/pkg/config"
	"github.com/stack-labs/stack-rpc/pkg/config/reader"
	cliSource "github.com/stack-labs/stack-rpc/pkg/config/source/cli"
	"github.com/stack-labs/stack-rpc/pkg/config/source/file"
	"github.com/stack-labs/stack-rpc/util/log"
)

type Broker struct {
	Address string `json:"address"`
	Name    string `json:"name" `
}

type Pool struct {
	Size json.Number `json:"size"`
	TTL  json.Number `json:"ttl"`
}

type ClientRequest struct {
	Retries json.Number `json:"retries"`
	Timeout json.Number `json:"timeout"`
}

type Client struct {
	Protocol string        `json:"protocol"`
	Pool     Pool          `json:"pool"`
	Request  ClientRequest `json:"request"`
}

type Registry struct {
	Address  string      `json:"address"`
	Interval json.Number `json:"interval"`
	Name     string      `json:"name"`
	TTL      json.Number `json:"ttl"`
}

type Metadata []string

func (m Metadata) Value(k string) string {
	for _, s := range m {
		kv := strings.Split(s, "=")
		if len(kv) == 2 && kv[0] == k {
			return kv[1]
		}
	}

	return ""
}

type Server struct {
	Address   string   `json:"address"`
	Advertise string   `json:"advertise"`
	ID        string   `json:"id"`
	Metadata  Metadata `json:"metadata"`
	Name      string   `json:"name"`
	Protocol  string   `json:"protocol"`
	Version   string   `json:"version"`
}

type Selector struct {
	Name string `json:"name"`
}

type Transport struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Logger struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}

type stack struct {
	Broker    Broker    `json:"broker"`
	Client    Client    `json:"client"`
	Profile   string    `json:"profile"`
	Registry  Registry  `json:"registry"`
	Runtime   string    `json:"runtime"`
	Server    Server    `json:"server"`
	Selector  Selector  `json:"selector"`
	Transport Transport `json:"transport"`
	Logger    Logger    `json:"logger"`
}

type Value struct {
	Stack stack `json:"stack"`
}

type Config interface {
	reader.Values
	Close() error
}

type stackConfig struct {
	config config.Config
}

// Init Stack's Config component
// Any developer Don't use this Func anywhere. Init works for Stack Framework only
func New(opts ...Option) (Config, error) {
	var o = Options{}
	for _, opt := range opts {
		opt(&o)
	}

	// need read from config file
	if len(o.FilePath) > 0 {
		log.Info("config read from file:", o.FilePath)
		o.Sources = append(o.Sources, file.NewSource(file.WithPath(o.FilePath)))
	}
	defaultSource, _ := json.Marshal(GetDefault())
	o.Sources = append(o.Sources,
		cliSource.NewSource(o.App, cliSource.Context(o.App.Context())),
		memory.NewSource(memory.WithJSON(defaultSource)),
	)

	c, err := config.NewConfig(config.Storage(true), config.Watch(false))
	if err != nil {
		return nil, err
	}
	if err := c.Load(o.Sources...); err != nil {
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
