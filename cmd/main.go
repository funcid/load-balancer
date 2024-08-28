package main

import (
	"flag"
	"load-balancer/internal"
	"load-balancer/internal/factory"
	"load-balancer/internal/service"
	"log"
	"net/http"
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

	balancerFactory := factory.NewBalancerFactory()
	balancer := balancerFactory.CreateBalancer(*balancerType, servers)
	lb := internal.NewLoadBalancer(balancer)

	go service.StartServer(8081)
	go service.StartServer(8082)
	go service.StartServer(8083)

	log.Println("Starting Load Balancer at :8080")
	http.HandleFunc("/", lb.HandleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
