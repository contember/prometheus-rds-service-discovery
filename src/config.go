package main

import (
	"fmt"
	"os"
	"strings"
)

type ScrapingTarget struct {
	Cluster string
	Region  string
}

type ScrapingConfig struct {
	Targets []ScrapingTarget
}

func buildFromEnvironments() (ScrapingConfig, error) {
	var scrapingTargets []ScrapingTarget

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		key, value := parts[0], parts[1]

		// Check if the environment variable matches the pattern SCRAPER_X_CLUSTER or SCRAPER_X_REGION
		if strings.HasPrefix(key, "SCRAPER_") && strings.HasSuffix(key, "_CLUSTER") {
			// Extract the X part from the variable name (SCRAPER_X_CLUSTER)
			suffix := strings.TrimSuffix(strings.TrimPrefix(key, "SCRAPER_"), "_CLUSTER")

			// Find the corresponding SCRAPER_X_REGION environment variable
			regionKey := fmt.Sprintf("SCRAPER_%s_REGION", suffix)
			regionValue := os.Getenv(regionKey)

			target := ScrapingTarget{
				Cluster: value,
				Region:  regionValue,
			}

			scrapingTargets = append(scrapingTargets, target)
		}
	}

	return ScrapingConfig{Targets: scrapingTargets}, nil
}
