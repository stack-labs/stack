package loader

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/stack-labs/stack-rpc/pkg/config/reader"
	"github.com/stack-labs/stack-rpc/pkg/config/source"
)

// Watcher lets you watch sources and returns a merged ChangeSet
type Watcher interface {
	// First call to next may return the current Snapshot
	// If you are watching a path then only the data from
	// that path is returned.
	Next() (*Snapshot, error)
	// Stop watching for changes
	Stop() error
}

type watcher struct {
	exit    chan bool
	path    []string
	value   reader.Value
	reader  reader.Reader
	updates chan reader.Value
}

func (w *watcher) Next() (*Snapshot, error) {
	for {
		select {
		case <-w.exit:
			return nil, errors.New("watcher stopped")
		case v := <-w.updates:
			if bytes.Equal(w.value.Bytes(), v.Bytes()) {
				continue
			}
			w.value = v

			cs := &source.ChangeSet{
				Data:      v.Bytes(),
				Format:    w.reader.String(),
				Source:    "memory",
				Timestamp: time.Now(),
			}
			cs.Sum()

			return &Snapshot{
				ChangeSet: cs,
				Version:   fmt.Sprintf("%d", time.Now().Unix()),
			}, nil
		}
	}
}

func (w *watcher) Stop() error {
	select {
	case <-w.exit:
	default:
		close(w.exit)
	}
	return nil
}
