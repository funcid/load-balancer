package algorithms

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"load-balancer/internal/balancing"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
)

type GeoServer struct {
	URL       *url.URL
	Latitude  float64
	Longitude float64
}

type GeographicalBalancer struct {
	servers []GeoServer
	mu      sync.Mutex
}

func NewGeographicalBalancer(servers []GeoServer) *GeographicalBalancer {
	return &GeographicalBalancer{servers: servers}
}

func (gb *GeographicalBalancer) GetCoordinates(ip string) (float64, float64, error) {
	apiKey := os.Getenv("GEO_API_KEY")
	if apiKey == "" {
		return 0, 0, errors.New("GEO_API_KEY environment variable not set")
	}

	resp, err := http.Get(fmt.Sprintf("https://api.ipgeolocation.io/ipgeo?apiKey=%s&ip=%s", apiKey, ip))
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, errors.New("failed to get geolocation info")
	}

	var result struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	if err = json.Unmarshal(body, &result); err != nil {
		return 0, 0, err
	}

	lat, err := strconv.ParseFloat(result.Latitude, 64)
	if err != nil {
		return 0, 0, err
	}

	lon, err := strconv.ParseFloat(result.Longitude, 64)
	if err != nil {
		return 0, 0, err
	}

	return lat, lon, nil
}

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in kilometers
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func (gb *GeographicalBalancer) NextBackend(r *http.Request) (*url.URL, error) {
	if len(gb.servers) == 0 {
		return nil, balancing.ErrNoAvailableBackends
	}

	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, err
	}

	clientLat, clientLon, err := gb.GetCoordinates(clientIP)
	if err != nil {
		return nil, err
	}

	gb.mu.Lock()
	defer gb.mu.Unlock()

	var closestServer *GeoServer
	shortestDistance := math.MaxFloat64

	for _, server := range gb.servers {
		dist := Haversine(clientLat, clientLon, server.Latitude, server.Longitude)
		if dist < shortestDistance {
			shortestDistance = dist
			closestServer = &server
		}
	}

	if closestServer == nil {
		return nil, balancing.ErrNoAvailableBackends
	}

	return closestServer.URL, nil
}
