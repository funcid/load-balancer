package balancing

import (
	"errors"
	"github.com/valyala/fasthttp"
	"net/url"
)

var ErrNoAvailableBackends = errors.New("no available backends")

type Balancer interface {
	NextBackend(ctx *fasthttp.RequestCtx) (*url.URL, error)
}

type ReleasableBalancer interface {
	Balancer
	ReleaseBackend(backend *url.URL)
}
