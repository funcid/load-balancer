package algorithms

import (
	"github.com/valyala/fasthttp"
	"load-balancer/internal/balancing"
	"net/url"
	"sync"
)

type LeastConnectionsBalancer struct {
	backends    []*url.URL
	connections []int
	mu          sync.Mutex
}

func NewLeastConnectionsBalancer(backends []*url.URL) *LeastConnectionsBalancer {
	return &LeastConnectionsBalancer{
		backends:    backends,
		connections: make([]int, len(backends)),
	}
}

func (lc *LeastConnectionsBalancer) NextBackend(*fasthttp.RequestCtx) (*url.URL, error) {
	if len(lc.backends) == 0 {
		return nil, balancing.ErrNoAvailableBackends
	}

	lc.mu.Lock()
	defer lc.mu.Unlock()

	minIndex := 0
	for i, count := range lc.connections {
		if count < lc.connections[minIndex] {
			minIndex = i
		}
	}

	lc.connections[minIndex]++
	return lc.backends[minIndex], nil
}

func (lc *LeastConnectionsBalancer) ReleaseBackend(backend *url.URL) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	for i, b := range lc.backends {
		if b.String() == backend.String() {
			lc.connections[i]--
			break
		}
	}
}
