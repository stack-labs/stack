package config

import (
	"encoding/json"
	"fmt"
	"github.com/stack-labs/stack-rpc/util/log"
	"strings"

	"github.com/stack-labs/stack-rpc/pkg/config"
	"github.com/stack-labs/stack-rpc/pkg/config/reader"
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

type Stack struct {
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
	Stack Stack `json:"stack"`
}

type Config interface {
	reader.Values
	Init(opts ...Option) error
	Close() error
}

type stackConfig struct {
	config config.Config
	opts   Options
}

// Init Stack's Config component
// Any developer Don't use this Func anywhere. NewConfig works for Stack Framework only
func NewConfig(opts ...Option) Config {
	var o = Options{
		Watch: true,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return &stackConfig{opts: o}
}

func (c *stackConfig) Init(opts ...Option) (err error) {
	for _, opt := range opts {
		opt(&c.opts)
	}

	defer func() {
		if err != nil {
			log.Errorf("config init error: %s", err)
		}
	}()

	cfg, err := config.NewConfig(
		config.Storage(c.opts.Storage),
		config.Watch(c.opts.Watch),
	)
	if err != nil {
		err = fmt.Errorf("create new config error: %s", err)
		return
	}

	if err = cfg.Load(c.opts.Sources...); err != nil {
		err = fmt.Errorf("load sources error: %s", err)
		return
	}

	c.config = cfg

	// cache c as sugar
	_sugar = c
	return nil
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
