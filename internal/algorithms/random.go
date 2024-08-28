package algorithms

import (
	"load-balancer/internal/balancing"
	"math/rand"
	"net/http"
	"net/url"
)

type RandomBalancer struct {
	backends []*url.URL
}

func NewRandomBalancer(backends []*url.URL) *RandomBalancer {
	return &RandomBalancer{backends: backends}
}

func (rb *RandomBalancer) NextBackend(*http.Request) (*url.URL, error) {
	if len(rb.backends) == 0 {
		return nil, balancing.ErrNoAvailableBackends
	}
	index := rand.Intn(len(rb.backends))
	return rb.backends[index], nil
}
