package balancing

import (
	"github.com/valyala/fasthttp"
)

type LoadBalancer struct {
	balancer Balancer
}

func NewLoadBalancer(balancer Balancer) *LoadBalancer {
	return &LoadBalancer{balancer: balancer}
}

func (lb *LoadBalancer) HandleRequest(ctx *fasthttp.RequestCtx) {
	backend, err := lb.balancer.NextBackend(ctx)
	if err != nil {
		ctx.Error("Service Unavailable", fasthttp.StatusServiceUnavailable)
		return
	}

	// Это URL бэкенда, к которому будем проксировать запрос
	backendURL := backend.String()

	// Создаем клиент для отправки запросов к бэкенду
	client := &fasthttp.HostClient{
		Addr: backend.Host,
	}

	if backend.Scheme == "https" {
		client.IsTLS = true
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	ctx.Request.CopyTo(req)
	req.SetRequestURI(backendURL)

	if err := client.Do(req, resp); err != nil {
		ctx.Error("Bad Gateway", fasthttp.StatusBadGateway)
		return
	}

	resp.CopyTo(&ctx.Response)
}
