package algorithms

import (
	"github.com/valyala/fasthttp"
	"load-balancer/internal/balancing"
	"net/url"
	"sync"
)

type LeastActiveBalancer struct {
	backends []*url.URL
	active   []int
	mu       sync.Mutex
}

func NewLeastActiveBalancer(backends []*url.URL) *LeastActiveBalancer {
	return &LeastActiveBalancer{
		backends: backends,
		active:   make([]int, len(backends)),
	}
}

func (la *LeastActiveBalancer) NextBackend(*fasthttp.RequestCtx) (*url.URL, error) {
	if len(la.backends) == 0 {
		return nil, balancing.ErrNoAvailableBackends
	}

	la.mu.Lock()
	defer la.mu.Unlock()

	minIndex := 0
	for i, count := range la.active {
		if count < la.active[minIndex] {
			minIndex = i
		}
	}

	la.active[minIndex]++
	return la.backends[minIndex], nil
}

func (la *LeastActiveBalancer) ReleaseBackend(backend *url.URL) {
	la.mu.Lock()
	defer la.mu.Unlock()

	for i, b := range la.backends {
		if b.String() == backend.String() {
			la.active[i]--
			break
		}
	}
}
