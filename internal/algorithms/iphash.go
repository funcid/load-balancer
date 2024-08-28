package algorithms

import (
	"github.com/valyala/fasthttp"
	"hash/fnv"
	"load-balancer/internal/balancing"
	"net/url"
)

type IpHashBalancer struct {
	backends []*url.URL
}

func NewIPHashBalancer(backends []*url.URL) *IpHashBalancer {
	return &IpHashBalancer{backends: backends}
}

func (ip *IpHashBalancer) NextBackend(ctx *fasthttp.RequestCtx) (*url.URL, error) {
	if len(ip.backends) == 0 {
		return nil, balancing.ErrNoAvailableBackends
	}
	clientIP := ctx.RemoteIP().String()
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(clientIP))
	index := hash.Sum32() % uint32(len(ip.backends))
	return ip.backends[index], nil
}
