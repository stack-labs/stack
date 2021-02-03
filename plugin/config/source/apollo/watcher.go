package apollo

import (
	"github.com/stack-labs/stack-rpc/pkg/config/source"
	"github.com/stack-labs/stack-rpc/plugin/config/source/apollo/agollo"
	"github.com/stack-labs/stack-rpc/plugin/config/source/apollo/agollo/storage"
)

type changeListener struct {
	change chan<- bool
}

func (l *changeListener) OnChange(event *storage.ChangeEvent) {
	l.change <- true
}

func (l *changeListener) OnNewestChange(event *storage.FullChangeEvent) {
	// ignore this change, OnChange will do every change thing.
	// l.change <- true
}

type watcher struct {
	client     *agollo.Client
	listener   *changeListener
	namespaces []string
	change     <-chan bool
}

func newWatcher(namespaces []string, client *agollo.Client) (*watcher, error) {
	change := make(chan bool)
	w := &watcher{
		client:     client,
		namespaces: namespaces,
		listener: &changeListener{
			change: change,
		},
		change: change,
	}

	w.client.AddChangeListener(w.listener)

	return w, nil
}

func (w *watcher) Next() (*source.ChangeSet, error) {
	select {
	case <-w.change:
		return read(w.namespaces, w.client)
	}
}

func (w *watcher) Stop() error {
	w.client.RemoveChangeListener(w.listener)
	return nil
}
