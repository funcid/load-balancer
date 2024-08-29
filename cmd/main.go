package main

import (
	"flag"
	"github.com/valyala/fasthttp"
	"load-balancer/internal/balancer_factory"
	"load-balancer/internal/balancing"
	"load-balancer/internal/service"
	"log"
)

func main() {
	balancerType := flag.String(
		"balancer",
		"roundrobin",
		"Balancer type: roundrobin, iphash, leastactive, geographical, leastconnections, random",
	)
	flag.Parse()

	servers := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	// Используем фабрику для создания балансировщика
	balancer := balancer_factory.CreateBalancer(*balancerType, servers)
	lb := balancing.NewLoadBalancer(balancer)

	// Запускаем серверы
	go service.StartServer(8081)
	go service.StartServer(8082)
	go service.StartServer(8083)

	log.Println("Starting Load Balancer at :8080")

	// Запускаем сервер балансировщика с использованием fasthttp
	log.Fatal(fasthttp.ListenAndServe(":8080", func(ctx *fasthttp.RequestCtx) {
		lb.HandleRequest(ctx)
	}))
}
