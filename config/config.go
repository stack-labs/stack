package config

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/stack-labs/stack-rpc/pkg/config"
	"github.com/stack-labs/stack-rpc/pkg/config/reader"
	"github.com/stack-labs/stack-rpc/util/log"
)

var (
	// Define the tag name for setting autowired value of Options
	// sc equals stack-config :)
	// todo support custom tagName
	DefaultOptionsTagName     = "sc"
	DefaultHierarchySeparator = "."

	// holds all the Options
	optionsPool = make(map[string]reflect.Value)
)

type Broker struct {
	Address string `json:"address" sc:"address"`
	Name    string `json:"name" sc:"name"`
}

type Pool struct {
	Size json.Number `json:"size" sc:"size"`
	TTL  json.Number `json:"ttl" sc:"ttl"`
}

type ClientRequest struct {
	Retries json.Number `json:"retries" sc:"retries"`
	Timeout json.Number `json:"timeout" sc:"timeout"`
}

type Client struct {
	Protocol string        `json:"protocol" sc:"protocol"`
	Pool     Pool          `json:"pool" sc:"pool"`
	Request  ClientRequest `json:"request" sc:"request"`
}

type Registry struct {
	Address  string      `json:"address" sc:"address"`
	Interval json.Number `json:"interval" sc:"interval"`
	Name     string      `json:"name" sc:"name"`
	TTL      json.Number `json:"ttl" sc:"ttl"`
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
	Address   string   `json:"address" sc:"address"`
	Advertise string   `json:"advertise" sc:"advertise"`
	ID        string   `json:"id" sc:"id"`
	Metadata  Metadata `json:"metadata" sc:"metadata"`
	Name      string   `json:"name" sc:"name"`
	Protocol  string   `json:"protocol" sc:"protocol"`
	Version   string   `json:"version" sc:"version"`
}

type Selector struct {
	Name string `json:"name" sc:"name"`
}

type Transport struct {
	Name    string `json:"name" sc:"name"`
	Address string `json:"address" sc:"address"`
}

type Logger struct {
	Name  string `json:"name" sc:"name"`
	Level string `json:"level" sc:"level"`
}

type Stack struct {
	Broker    Broker    `json:"broker" sc:"broker"`
	Client    Client    `json:"client" sc:"client"`
	Profile   string    `json:"profile" sc:"profile"`
	Registry  Registry  `json:"registry" sc:"registry"`
	Runtime   string    `json:"runtime" sc:"runtime"`
	Server    Server    `json:"server" sc:"server"`
	Selector  Selector  `json:"selector" sc:"selector"`
	Transport Transport `json:"transport" sc:"transport"`
	Logger    Logger    `json:"logger" sc:"logger"`
}

type Value struct {
	Stack Stack `json:"stack" sc:"stack"`
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

func (c *stackConfig) Init(opts ...Option) (err error) {
	for _, opt := range opts {
		opt(&c.opts)
	}

	if c.opts.Context == nil {
		c.opts.Context = context.Background()
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
	// set the autowired values
	injectAutowired(c.opts.Context)

	return nil
}

func (c *stackConfig) Get(path ...string) reader.Value {
	tempPath := path
	if len(path) == 1 {
		if strings.Contains(path[0], DefaultHierarchySeparator) {
			tempPath = strings.Split(path[0], DefaultHierarchySeparator)
		}
	}

	return c.config.Get(tempPath...)
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

func RegisterOptions(options ...interface{}) {
	for _, option := range options {
		val := reflect.ValueOf(option)
		if val.Kind() != reflect.Ptr {
			log.Error("options must be a pointer")
			return
		}

		_, file, line, _ := runtime.Caller(0)

		key := fmt.Sprintf("%s%d", file, line)

		optionsPool[key] = val
	}
}
