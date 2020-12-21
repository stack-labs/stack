package config

import (
	"context"
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

		_, file, line, _ := runtime.Caller(1)

		key := fmt.Sprintf("%s#L%d", file, line)

		optionsPool[key] = val
	}
}
