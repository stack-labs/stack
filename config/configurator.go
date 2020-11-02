package config

import (
	"fmt"
	"sync"

	log "github.com/stack-labs/stack-rpc/logger"
)

var (
	m      sync.RWMutex
	inited bool

	c = &configurator{}
)

// Configurator
type Configurator interface {
	App(name string, config interface{}) (err error)
}

// configurator
type configurator struct {
	conf Config
}

func (c *configurator) App(name string, config interface{}) (err error) {
	v := c.conf.Get(name)
	if v != nil {
		err = v.Scan(config)
	} else {
		err = fmt.Errorf("[App] 配置不存在，err：%s", name)
	}

	return
}

// C returns the Configurator
func C() Configurator {
	return c
}

func (c *configurator) init(ops Options) (err error) {
	m.Lock()
	defer m.Unlock()

	if inited {
		log.Info("[init] 配置已经初始化过")
		return
	}

	c.conf, err = NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	// 加载配置
	err = c.conf.Load(ops.Source...)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Info("[init] 侦听配置变动 ...")

		// 开始侦听变动事件
		watcher, err := c.conf.Watch()
		if err != nil {
			log.Fatal(err)
		}

		for {
			v, err := watcher.Next()
			if err != nil {
				log.Fatal(err)
			}

			log.Infof("[init] 侦听配置变动: %v", string(v.Bytes()))
		}
	}()

	// 标记已经初始化
	inited = true
	return
}

// Init 初始化配置
func Init(opts ...Option) error {
	ops := Options{}
	for _, o := range opts {
		o(&ops)
	}

	c = &configurator{}

	return c.init(ops)
}
