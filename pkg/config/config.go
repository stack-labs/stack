// Package config is an interface for dynamic configuration.
package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/pkg/config/loader"
	"github.com/stack-labs/stack-rpc/pkg/config/reader"
	"github.com/stack-labs/stack-rpc/pkg/config/source"
	"github.com/stack-labs/stack-rpc/pkg/config/storage"
	"github.com/stack-labs/stack-rpc/pkg/config/storage/file"
)

// Config is an interface abstraction for dynamic configuration
type Config interface {
	// provide the reader.Values interface
	reader.Values
	// Stop the config loader/watcher
	Close() error
	// Load config sources
	Load(source ...source.Source) error
	// Force a source change set sync
	Sync() error
	// Watch a value for changes
	Watch(path ...string) (Watcher, error)
}

type config struct {
	exit    chan bool
	storage storage.Storage
	loader  loader.Loader
	opts    Options

	sync.RWMutex
	// the current snapshot
	snap *loader.Snapshot
	// the current values
	values reader.Values
}

// NewConfig returns new config
func NewConfig(opts ...Option) (Config, error) {
	return newConfig(opts...)
}

func newConfig(opts ...Option) (Config, error) {
	options := NewOptions(opts...)

	l := loader.NewLoader(loader.WithWatch(options.Watch))
	if err := l.Load(); err != nil {
		return nil, err
	}
	snap, err := l.Snapshot()
	if err != nil {
		return nil, err
	}
	values, err := l.Values(snap.ChangeSet)
	if err != nil {
		return nil, err
	}

	var cStorage storage.Storage
	if options.Storage {
		dir := options.StorageDir
		if len(dir) == 0 {
			dir, err = os.Getwd()
			if err != nil {
				return nil, err
			}
		}
		f := filepath.Join(dir, "stack_config.conf")
		cStorage = file.NewStorage(f)
	}

	c := &config{
		exit:    make(chan bool),
		storage: cStorage,
		opts:    options,
		snap:    snap,
		loader:  l,
		values:  values,
	}

	if c.opts.Watch {
		go c.run()
	}

	return c, nil
}

func (c *config) writeStorage(snap *loader.Snapshot) {
	if snap != nil && c.opts.Storage && c.storage != nil {
		var out bytes.Buffer
		var bs []byte
		// todo support more types
		// beautify here is not a good option (:.
		err := json.Indent(&out, snap.ChangeSet.Data, "", "  ")
		if err != nil {
			log.Errorf("beautify data error: %v", err)
			bs = snap.ChangeSet.Data
		} else {
			bs = out.Bytes()
		}

		if err := c.storage.Write(bs); err != nil {
			log.Errorf("config storage write error: %v", err)
		}
	}
}

func (c *config) run() {
	watch := func(w loader.Watcher) error {
		for {
			// get change set
			snap, err := w.Next()
			if err != nil {
				return err
			}

			c.writeStorage(snap)

			c.Lock()

			c.snap = snap
			c.values, _ = c.loader.Values(snap.ChangeSet)

			c.Unlock()
		}
	}

	for {
		w, err := c.loader.Watch()
		if err != nil {
			log.Warnf("create loader watcher error: %v", err)
			time.Sleep(time.Second)
			continue
		}

		done := make(chan bool)

		// the stop watch func
		go func() {
			select {
			case <-done:
			case <-c.exit:
			}
			_ = w.Stop()
		}()

		// block watch
		if err := watch(w); err != nil {
			// do something better
			time.Sleep(time.Second)
		}

		// close done chan
		close(done)

		// if the config is closed exit
		select {
		case <-c.exit:
			return
		default:
		}
	}
}

func (c *config) Map() map[string]interface{} {
	c.RLock()
	defer c.RUnlock()
	return c.values.Map()
}

func (c *config) Scan(v interface{}) error {
	c.RLock()
	defer c.RUnlock()
	return c.values.Scan(v)
}

// sync loads all the sources, calls the parser and updates the config
func (c *config) Sync() error {
	if err := c.loader.Sync(); err != nil {
		return err
	}

	snap, err := c.loader.Snapshot()
	if err != nil {
		return err
	}

	c.writeStorage(snap)

	c.Lock()
	defer c.Unlock()

	c.snap = snap
	values, err := c.loader.Values(snap.ChangeSet)
	if err != nil {
		return err
	}
	c.values = values

	return nil
}

func (c *config) Close() error {
	select {
	case <-c.exit:
		return nil
	default:
		close(c.exit)
		if err := c.loader.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (c *config) Get(path ...string) reader.Value {
	c.RLock()
	defer c.RUnlock()

	// did sync actually work?
	if c.values != nil {
		return c.values.Get(path...)
	}

	// no value
	return newValue()
}

func (c *config) Bytes() []byte {
	c.RLock()
	defer c.RUnlock()

	if c.values == nil {
		return []byte{}
	}

	return c.values.Bytes()
}

func (c *config) loadBackupConfig() error {
	bytes, err := c.storage.Load()
	if err != nil {
		return err
	}

	cs := &source.ChangeSet{
		Data:      bytes,
		Format:    "json", // only json reader
		Source:    "backup",
		Timestamp: time.Now(),
	}
	cs.Sum()
	snap := &loader.Snapshot{
		ChangeSet: cs,
		Version:   fmt.Sprintf("%d", time.Now().Unix()),
	}

	c.Lock()
	defer c.Unlock()

	c.snap = snap
	values, err := c.loader.Values(snap.ChangeSet)
	if err != nil {
		return err
	}
	c.values = values

	return nil
}

func (c *config) Load(sources ...source.Source) error {
	if err := c.loader.Load(sources...); err != nil {
		if c.opts.Storage && c.storage.Exist() {
			log.Warn("load config from backup file")
			return c.loadBackupConfig()
		}
		return err
	}

	snap, err := c.loader.Snapshot()
	if err != nil {
		return err
	}
	c.writeStorage(snap)

	c.Lock()
	defer c.Unlock()

	c.snap = snap
	values, err := c.loader.Values(snap.ChangeSet)
	if err != nil {
		return err
	}
	c.values = values

	return nil
}

func (c *config) Watch(path ...string) (Watcher, error) {
	value := c.Get(path...)

	w, err := c.loader.Watch(path...)
	if err != nil {
		return nil, err
	}

	return &watcher{
		lw:    w,
		rd:    c.loader.Reader(),
		path:  path,
		value: value,
	}, nil
}

func (c *config) String() string {
	return "config"
}
