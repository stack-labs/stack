package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	log "github.com/stack-labs/stack-rpc/logger"
)

var (
	m      sync.RWMutex
	c      *configurator
	inited bool
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
		err = fmt.Errorf("config isn't existed，err：%s", name)
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
		log.Info("config has been inited")
		return
	}

	c.conf, err = NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// set the static dir to working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Errorf("set log static dir error: %s", err)
		return
	}

	DefaultStaticDir = dir + string(filepath.Separator) + staticDirName

	// 加载配置
	err = c.conf.Load(ops.Source...)
	if err != nil {
		log.Fatal(err)
	}

	go c.watchChanges()

	inited = true
	return
}

func (c *configurator) watchChanges() {
	log.Infof("start to watch config changes")
	watcher, err := c.conf.Watch()
	if err != nil {
		log.Errorf("watch config changes error: %s", err)
	}

	for {
		v, err := watcher.Next()
		if err != nil {
			log.Errorf("watcher get next change error: %s", err)
		}

		log.Debugf("configuration changes: %v", string(v.Bytes()))
	}
}

func Init(opts ...Option) error {
	ops := Options{}
	for _, o := range opts {
		o(&ops)
	}

	c = &configurator{}

	return c.init(ops)
}
