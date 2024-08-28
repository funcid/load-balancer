package algorithms

import (
	"github.com/valyala/fasthttp"
	"load-balancer/internal/balancing"
	"math/rand"
	"net/url"
)

type RandomBalancer struct {
	backends []*url.URL
}

func NewRandomBalancer(backends []*url.URL) *RandomBalancer {
	return &RandomBalancer{backends: backends}
}

func (rb *RandomBalancer) NextBackend(*fasthttp.RequestCtx) (*url.URL, error) {
	if len(rb.backends) == 0 {
		return nil, balancing.ErrNoAvailableBackends
	}
	index := rand.Intn(len(rb.backends))
	return rb.backends[index], nil
}
