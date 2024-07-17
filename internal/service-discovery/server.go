package service_discovery

import (
	"context"
	"encoding/json"
	service_scraper "github.com/contember/cloud-prometheus-rds-service-discovery/internal/service-scraper"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type ServiceDiscovery struct {
	serverAddress string
	scraper       *service_scraper.Scraper
}

func NewServiceDiscovery(address string, scraper *service_scraper.Scraper) *ServiceDiscovery {
	return &ServiceDiscovery{
		serverAddress: address,
		scraper:       scraper,
	}
}

func (s *ServiceDiscovery) Serve(_ context.Context, errors chan error) {
	http.HandleFunc("/discovery", s.serverGetClusters)
	http.HandleFunc("/.health", s.serverGetHealth)

	log.Println("RDS Instances service discovery running on ", s.serverAddress)
	log.Println("Flowing cluster ready for scraping:", s.scraper.Targets)

	errors <- http.ListenAndServe(s.serverAddress, nil)
}

// @see https://prometheus.io/docs/prometheus/latest/http_sd/#writing-http-service-discovery
func (s *ServiceDiscovery) serverGetClusters(w http.ResponseWriter, r *http.Request) {
	targets, err := s.scraper.GetInstances(context.TODO())
	if err != nil {
		log.Errorf("Error during scraping RDS instances %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(targets)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the content type and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ServiceDiscovery) serverGetHealth(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Service running!")
}
