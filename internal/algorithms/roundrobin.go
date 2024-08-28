package algorithms

import (
	"github.com/valyala/fasthttp"
	"load-balancer/internal/balancing"
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

func (rr *RoundRobinBalancer) NextBackend(*fasthttp.RequestCtx) (*url.URL, error) {
	if len(rr.backends) == 0 {
		return nil, balancing.ErrNoAvailableBackends
	}
	next := atomic.AddUint64(&rr.count, 1)
	return rr.backends[next%uint64(len(rr.backends))], nil
}
