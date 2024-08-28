package balancing

import (
	"net/http"
	"net/http/httputil"
)

type LoadBalancer struct {
	balancer Balancer
}

func NewLoadBalancer(balancer Balancer) *LoadBalancer {
	return &LoadBalancer{balancer: balancer}
}

func (lb *LoadBalancer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	backend, err := lb.balancer.NextBackend(r)
	if err != nil {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(backend)

	if releasableBalancer, ok := lb.balancer.(ReleasableBalancer); ok {
		// Обработчик ошибок прокси-сервера для уменьшения счетчика
		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
			releasableBalancer.ReleaseBackend(backend)
			http.Error(rw, "Bad Gateway", http.StatusBadGateway)
		}

		r.URL.Host = backend.Host
		r.URL.Scheme = backend.Scheme
		r.Host = backend.Host

		proxy.ServeHTTP(w, r)

		// Освобождаем соединение после успешного завершения
		releasableBalancer.ReleaseBackend(backend)
	} else {
		r.URL.Host = backend.Host
		r.URL.Scheme = backend.Scheme
		r.Host = backend.Host

		proxy.ServeHTTP(w, r)
	}
}
