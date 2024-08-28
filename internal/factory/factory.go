package factory

import (
	"load-balancer/internal"
	"load-balancer/internal/balancers"
	"log"
	"net/url"
)

type BalancerFactory struct{}

func NewBalancerFactory() *BalancerFactory {
	return &BalancerFactory{}
}

func (f *BalancerFactory) CreateBalancer(balancerType string, servers []string) internal.Balancer {
	var backends []*url.URL
	for _, server := range servers {
		u, err := url.Parse(server)
		if err != nil {
			log.Fatalf("Error parsing server URL: %v", err)
		}
		backends = append(backends, u)
	}

	var balancer internal.Balancer

	switch balancerType {
	case "roundrobin":
		balancer = balancers.NewRoundRobinBalancer(backends)
	case "iphash":
		balancer = balancers.NewIPHashBalancer(backends)
	case "leastactive":
		balancer = balancers.NewLeastActiveBalancer(backends)
	case "leastconnections":
		balancer = balancers.NewLeastConnectionsBalancer(backends)
	case "random":
		balancer = balancers.NewRandomBalancer(backends)
	case "geographical":
		balancer = balancers.NewGeographicalBalancer(
			[]balancers.GeoServer{
				{URL: backends[0], Latitude: 37.7749, Longitude: -122.4194}, // San Francisco
				{URL: backends[1], Latitude: 40.7128, Longitude: -74.0060},  // New York
				{URL: backends[2], Latitude: 51.5074, Longitude: -0.1278},   // London
			},
		)
	default:
		log.Fatalf("Unknown balancer type: %s", balancerType)
	}

	return balancer
}
