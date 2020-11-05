// Package config is an interface for dynamic configuration.
package config

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/stack-labs/stack-rpc/config/loader"
	"github.com/stack-labs/stack-rpc/config/reader"
	"github.com/stack-labs/stack-rpc/config/source"
	"github.com/stack-labs/stack-rpc/config/storage"
	"github.com/stack-labs/stack-rpc/config/storage/file"
	log "github.com/stack-labs/stack-rpc/logger"
)

// Config is an interface abstraction for dynamic configuration
type Config interface {
	// provide the reader.Values interface
	reader.Values
	// Stop the config loader/watcher
	Close() error
	// Load config sources
	Load(source ...source.Source) error
	// Force a source changeset sync
	Sync() error
	// Watch a value for changes
	Watch(path ...string) (Watcher, error)
}

type config struct {
	exit    chan bool
	storage storage.Storage
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

	if err := options.Loader.Load(); err != nil {
		return nil, err
	}
	snap, err := options.Loader.Snapshot()
	if err != nil {
		return nil, err
	}
	values, err := options.Reader.Values(snap.ChangeSet)
	if err != nil {
		return nil, err
	}

	var cStorage storage.Storage
	if options.EnableStorage {
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		f := fmt.Sprintf("%s/.stack_config/config.conf", dir)
		cStorage = file.NewStorage(f)
	}

	c := &config{
		exit:    make(chan bool),
		storage: cStorage,
		opts:    options,
		snap:    snap,
		values:  values,
	}

	go c.run()

	return c, nil
}

func (c *config) writeStorage(snap *loader.Snapshot) {
	if snap != nil && c.opts.EnableStorage && c.storage != nil {
		if err := c.storage.Write(snap.ChangeSet.Data); err != nil {
			log.Errorf("config storage write error: %v", err)
		}
	}
}

func (c *config) run() {
	watch := func(w loader.Watcher) error {
		for {
			// get changeset
			snap, err := w.Next()
			if err != nil {
				return err
			}

			c.writeStorage(snap)

			c.Lock()

			c.snap = snap
			c.values, _ = c.opts.Reader.Values(snap.ChangeSet)

			c.Unlock()
		}
	}

	for {
		w, err := c.opts.Loader.Watch()
		if err != nil {
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
			w.Stop()
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
	if err := c.opts.Loader.Sync(); err != nil {
		return err
	}

	snap, err := c.opts.Loader.Snapshot()
	if err != nil {
		return err
	}

	c.writeStorage(snap)

	c.Lock()
	defer c.Unlock()

	c.snap = snap
	vals, err := c.opts.Reader.Values(snap.ChangeSet)
	if err != nil {
		return err
	}
	c.values = vals

	return nil
}

func (c *config) Close() error {
	select {
	case <-c.exit:
		return nil
	default:
		close(c.exit)
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
		Format:    c.opts.Reader.String(),
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
	values, err := c.opts.Reader.Values(snap.ChangeSet)
	if err != nil {
		return err
	}
	c.values = values

	return nil
}

func (c *config) Load(sources ...source.Source) error {
	if err := c.opts.Loader.Load(sources...); err != nil {
		if c.opts.EnableStorage && c.storage.Exist() {
			log.Warn("load config from backup file")
			return c.loadBackupConfig()
		}
		return err
	}

	snap, err := c.opts.Loader.Snapshot()
	if err != nil {
		return err
	}
	c.writeStorage(snap)

	c.Lock()
	defer c.Unlock()

	c.snap = snap
	values, err := c.opts.Reader.Values(snap.ChangeSet)
	if err != nil {
		return err
	}
	c.values = values

	return nil
}

func (c *config) Watch(path ...string) (Watcher, error) {
	value := c.Get(path...)

	w, err := c.opts.Loader.Watch(path...)
	if err != nil {
		return nil, err
	}

	return &watcher{
		lw:    w,
		rd:    c.opts.Reader,
		path:  path,
		value: value,
	}, nil
}

func (c *config) String() string {
	return "config"
}
