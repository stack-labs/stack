// package loader manages loading from multiple sources
package loader

import (
	"container/list"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/pkg/config/reader"
	"github.com/stack-labs/stack-rpc/pkg/config/source"
)

// Loader manages loading sources
type Loader interface {
	// Stop the loader
	Close() error
	// Load the sources
	Load(...source.Source) error
	// A Snapshot of loaded config
	Snapshot() (*Snapshot, error)
	// Force sync of sources
	Sync() error
	// Watch for changes
	Watch(...string) (Watcher, error)

	Values(*source.ChangeSet) (reader.Values, error)
	Reader() reader.Reader
}

type loader struct {
	exit chan bool
	opts Options

	sync.RWMutex
	// the current snapshot
	snap *Snapshot
	// the current values
	values reader.Values
	// all the sets
	sets []*source.ChangeSet
	// all the sources
	sources []source.Source

	watchers *list.List
}

func NewLoader(opts ...Option) Loader {
	options := NewOptions(opts...)

	return &loader{
		exit:     make(chan bool),
		opts:     options,
		watchers: list.New(),
	}
}

func (m *loader) Values(c *source.ChangeSet) (reader.Values, error) {
	return m.opts.Reader.Values(c)
}

func (m *loader) Reader() reader.Reader {
	return m.opts.Reader
}

// Snapshot returns a snapshot of the current loaded config
func (m *loader) Snapshot() (*Snapshot, error) {
	if !m.loaded() {
		// not loaded, sync
		if err := m.Sync(); err != nil {
			return nil, err
		}
	}

	m.RLock()
	snap := m.snap.Clone()
	m.RUnlock()

	return snap, nil
}

// Sync loads all the sources, calls the parser and updates the config
func (m *loader) Sync() error {
	//nolint:prealloc
	var sets []*source.ChangeSet

	m.Lock()

	// read the source
	var gerr []string

	for _, source := range m.sources {
		ch, err := source.Read()
		if err != nil {
			gerr = append(gerr, err.Error())
			continue
		}
		sets = append(sets, ch)
	}

	// merge sets
	set, err := m.opts.Reader.Merge(sets...)
	if err != nil {
		m.Unlock()
		return err
	}

	// set values
	values, err := m.opts.Reader.Values(set)
	if err != nil {
		m.Unlock()
		return err
	}
	m.values = values
	m.snap = &Snapshot{
		ChangeSet: set,
		Version:   fmt.Sprintf("%d", time.Now().Unix()),
	}

	m.Unlock()

	// update watchers
	m.update()

	if len(gerr) > 0 {
		return fmt.Errorf("source loading errors: %s", strings.Join(gerr, "\n"))
	}

	return nil
}

func (m *loader) Close() error {
	select {
	case <-m.exit:
		return nil
	default:
		close(m.exit)
	}
	return nil
}

func (m *loader) Get(path ...string) (reader.Value, error) {
	if !m.loaded() {
		if err := m.Sync(); err != nil {
			return nil, err
		}
	}

	m.Lock()
	defer m.Unlock()

	// did sync actually work?
	if m.values != nil {
		return m.values.Get(path...), nil
	}

	// assuming vals is nil
	// create new vals

	ch := m.snap.ChangeSet

	// we are truly screwed, trying to load in a hacked way
	v, err := m.opts.Reader.Values(ch)
	if err != nil {
		return nil, err
	}

	// lets set it just because
	m.values = v

	if m.values != nil {
		return m.values.Get(path...), nil
	}

	// ok we're going hardcore now
	return nil, errors.New("no values")
}

func (m *loader) Load(sources ...source.Source) error {
	var gerrors []string

	for _, source := range sources {
		set, err := source.Read()
		if err != nil {
			gerrors = append(gerrors,
				fmt.Sprintf("error loading source %s: %v",
					source,
					err))
			// continue processing
			continue
		}
		m.Lock()
		m.sources = append(m.sources, source)
		m.sets = append(m.sets, set)
		idx := len(m.sets) - 1
		m.Unlock()
		if m.opts.Watch {
			go m.watch(idx, source)
		}
	}

	if err := m.reload(); err != nil {
		gerrors = append(gerrors, err.Error())
	}

	// Return errors
	if len(gerrors) != 0 {
		return errors.New(strings.Join(gerrors, "\n"))
	}
	return nil
}

func (m *loader) Watch(path ...string) (Watcher, error) {
	value, err := m.Get(path...)
	if err != nil {
		return nil, err
	}

	m.Lock()

	w := &watcher{
		exit:    make(chan bool),
		path:    path,
		value:   value,
		reader:  m.opts.Reader,
		updates: make(chan reader.Value, 1),
	}

	e := m.watchers.PushBack(w)

	m.Unlock()

	go func() {
		<-w.exit
		m.Lock()
		m.watchers.Remove(e)
		m.Unlock()
	}()

	return w, nil
}

func (m *loader) watch(idx int, s source.Source) {
	// watches a source for changes
	watch := func(idx int, s source.Watcher) error {
		for {
			// get changeset
			cs, err := s.Next()
			if err != nil {
				return err
			}

			m.Lock()
			m.sets[idx] = cs
			m.Unlock()

			if err := m.reload(); err != nil {
				return err
			}
		}
	}

	for {
		// watch the source
		w, err := s.Watch()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		done := make(chan bool)

		// the stop watch func
		go func() {
			select {
			case <-done:
			case <-m.exit:
			}
			_ = w.Stop()
		}()

		// block watch
		if err := watch(idx, w); err != nil {
			if err != source.ErrWatcherStopped {
				log.Errorf("loader watch source error : %s", err.Error())
			}
			// do something better
			time.Sleep(time.Second)
		}

		// close done chan
		close(done)

		// if the config is closed exit
		select {
		case <-m.exit:
			return
		default:
		}
	}
}

func (m *loader) loaded() bool {
	m.RLock()
	loaded := m.values != nil
	m.RUnlock()
	return loaded
}

// reload reads the sets and creates new values
func (m *loader) reload() error {
	m.Lock()

	// merge sets
	set, err := m.opts.Reader.Merge(m.sets...)
	if err != nil {
		m.Unlock()
		return err
	}

	// set values
	m.values, _ = m.opts.Reader.Values(set)
	m.snap = &Snapshot{
		ChangeSet: set,
		Version:   fmt.Sprintf("%d", time.Now().Unix()),
	}

	m.Unlock()

	// update watchers
	m.update()

	return nil
}

func (m *loader) update() {
	watchers := make([]*watcher, 0, m.watchers.Len())

	m.RLock()
	for e := m.watchers.Front(); e != nil; e = e.Next() {
		watchers = append(watchers, e.Value.(*watcher))
	}
	m.RUnlock()

	for _, w := range watchers {
		select {
		case w.updates <- m.values.Get(w.path...):
		default:
		}
	}
}
