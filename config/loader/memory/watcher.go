package memory

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/stack-labs/stack-rpc/config/loader"
	"github.com/stack-labs/stack-rpc/config/reader"
	"github.com/stack-labs/stack-rpc/config/source"
)

type watcher struct {
	exit    chan bool
	path    []string
	value   reader.Value
	reader  reader.Reader
	updates chan reader.Value
}

func (w *watcher) Next() (*loader.Snapshot, error) {
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

			return &loader.Snapshot{
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
