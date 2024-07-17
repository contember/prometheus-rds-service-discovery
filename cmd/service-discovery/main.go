package main

import (
	"context"
	service_discovery "github.com/contember/cloud-prometheus-rds-service-discovery/internal/service-discovery"
	service_scraper "github.com/contember/cloud-prometheus-rds-service-discovery/internal/service-scraper"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	signalChannel := make(chan os.Signal, 1)
	errChannel := make(chan error, 1)

	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	scraper, err := service_scraper.NewScraperFromEnvs()
	if err != nil {
		log.Fatalf("Error during setting up scraper via ENVS %s", err)
	}

	serviceDiscovery := service_discovery.NewServiceDiscovery(":"+port, scraper)
	go serviceDiscovery.Serve(context.Background(), errChannel)

	select {
	case err := <-errChannel:
		log.Errorf("Received error: %v", err)
		os.Exit(1)
	case sig := <-signalChannel:
		log.Infof("Received cancel signal: %v", sig)
		os.Exit(0)
	}
}
