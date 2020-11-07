package store

import (
	"github.com/stack-labs/stack-rpc/util/options"
)

// Nodes is a list of nodes used to back the store
func Nodes(a ...string) options.Option {
	return options.WithValue("store.nodes", a)
}

// Prefix sets a prefix to any key ids used
func Prefix(p string) options.Option {
	return options.WithValue("store.prefix", p)
}

// Namespace offers a way to have multiple isolated
// stores in the same backend, if supported.
func Namespace(n string) options.Option {
	return options.WithValue("store.namespace", n)
}
