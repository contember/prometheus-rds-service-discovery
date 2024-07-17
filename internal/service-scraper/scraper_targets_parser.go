package service_scraper

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ScraperTargetsEnvParser extracts AWS RDS cluster targets ARNs from environment variables
func scraperTargetsEnvParser() ([]ScrapingTarget, error) {
	var scrapingTargets []ScrapingTarget
	var scrapingEnvsRegexp = regexp.MustCompile(`^SCRAPER_(\d+)_CLUSTER_ARN$`)

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		key, value := parts[0], parts[1]
		keyMatches := scrapingEnvsRegexp.FindStringSubmatch(key)

		if len(keyMatches) > 1 {
			keyNumberStr := keyMatches[1]
			keyNumber, err := strconv.Atoi(keyNumberStr)
			if err != nil {
				return nil, fmt.Errorf("invalid key number: %s", keyNumberStr)
			}

			target := ScrapingTarget{
				ClusterArn:    value,
				ClusterNumber: keyNumber,
			}

			scrapingTargets = append(scrapingTargets, target)
		}
	}

	return scrapingTargets, nil
}
