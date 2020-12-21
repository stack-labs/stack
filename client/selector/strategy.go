package selector

import (
	"math/rand"
	"sync"
	"time"

	"github.com/stack-labs/stack-rpc/registry"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Strategy is a selection strategy e.g random, round robin
type Strategy func([]*registry.Service) (*registry.Node, error)

// Random is a random strategy algorithm for node selection
func Random() Strategy {
	return func(services []*registry.Service) (*registry.Node, error) {
		nodes := make([]*registry.Node, 0, len(services))

		for _, service := range services {
			nodes = append(nodes, service.Nodes...)
		}

		if len(nodes) == 0 {
			return nil, ErrNoneAvailable
		}

		i := rand.Int() % len(nodes)
		return nodes[i], nil
	}
}

// RoundRobin is a round robin strategy algorithm for node selection
func RoundRobin() Strategy {
	var i = rand.Int()
	var mtx sync.Mutex

	return func(services []*registry.Service) (*registry.Node, error) {
		nodes := make([]*registry.Node, 0, len(services))

		for _, service := range services {
			nodes = append(nodes, service.Nodes...)
		}

		if len(nodes) == 0 {
			return nil, ErrNoneAvailable
		}

		mtx.Lock()
		i++
		mtx.Unlock()

		return nodes[i%len(nodes)], nil
	}
}
