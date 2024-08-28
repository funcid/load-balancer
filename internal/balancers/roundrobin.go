package balancers

import (
	"load-balancer/internal/service"
	"net/http"
	"net/url"
	"sync/atomic"
)

type RoundRobinBalancer struct {
	backends []*url.URL
	count    uint64
}

func NewRoundRobinBalancer(backends []*url.URL) *RoundRobinBalancer {
	return &RoundRobinBalancer{backends: backends}
}

func (rr *RoundRobinBalancer) NextBackend(*http.Request) (*url.URL, error) {
	if len(rr.backends) == 0 {
		return nil, service.ErrNoAvailableBackends
	}
	next := atomic.AddUint64(&rr.count, 1)
	return rr.backends[next%uint64(len(rr.backends))], nil
}
