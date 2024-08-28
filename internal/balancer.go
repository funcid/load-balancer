package internal

import (
	"errors"
	"net/http"
	"net/url"
)

var ErrNoAvailableBackends = errors.New("no available backends")

type Balancer interface {
	NextBackend(r *http.Request) (*url.URL, error)
}

type ReleasableBalancer interface {
	Balancer
	ReleaseBackend(backend *url.URL)
}
