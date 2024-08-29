package balancer_factory

import (
	"load-balancer/internal/algorithms"
	"load-balancer/internal/balancing"
	"log"
	"net/url"
)

func CreateBalancer(balancerType string, servers []string) balancing.Balancer {
	var backends []*url.URL
	for _, server := range servers {
		u, err := url.Parse(server)
		if err != nil {
			log.Fatalf("Error parsing server URL: %v", err)
		}
		backends = append(backends, u)
	}

	var balancer balancing.Balancer

	switch balancerType {
	case "roundrobin":
		balancer = algorithms.NewRoundRobinBalancer(backends)
	case "iphash":
		balancer = algorithms.NewIPHashBalancer(backends)
	case "leastactive":
		balancer = algorithms.NewLeastActiveBalancer(backends)
	case "leastconnections":
		balancer = algorithms.NewLeastConnectionsBalancer(backends)
	case "random":
		balancer = algorithms.NewRandomBalancer(backends)
	case "geographical":
		balancer = algorithms.NewGeographicalBalancer(
			[]algorithms.GeoServer{
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
