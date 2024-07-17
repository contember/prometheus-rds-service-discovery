package service_scraper

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	log "github.com/sirupsen/logrus"
)

type ScrapingTarget struct {
	ClusterArn    string
	ClusterNumber int
}

type Scraper struct {
	Targets []ScrapingTarget
}

type ScraperDiscoveredTarget struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func NewScraperFromEnvs() (*Scraper, error) {
	targets, err := scraperTargetsEnvParser()
	if err != nil {
		return nil, err
	}

	return &Scraper{
		Targets: targets,
	}, nil
}

func NewScraper(targets []ScrapingTarget) *Scraper {
	return &Scraper{
		Targets: targets,
	}
}

func (s *Scraper) GetInstances(ctx context.Context) ([]ScraperDiscoveredTarget, error) {
	awsRdsInstances := make([]ScraperDiscoveredTarget, 0)
	awsSession, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	awsRds := rds.NewFromConfig(awsSession)

	// iterate over all registered clusters
	for _, cluster := range s.Targets {
		// filters cluster that contains cluster id that we want
		filters := []types.Filter{
			{
				Name:   aws.String("db-cluster-id"),
				Values: []string{cluster.ClusterArn},
			},
		}

		resp, err := awsRds.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{Filters: filters})
		if err != nil {
			log.Errorf("Error describing RDS cluster <%s> %v", cluster.ClusterArn, err)
			continue
		}

		// iterate over all instances in the cluster
		for _, instance := range resp.DBInstances {
			var instanceDsn = *instance.Endpoint.Address + ":" + fmt.Sprint(*instance.Endpoint.Port)
			var instanceTarget = ScraperDiscoveredTarget{
				Targets: []string{instanceDsn},
				Labels: map[string]string{
					"cluster_arn":         cluster.ClusterArn,
					"cluster_identifier":  *instance.DBClusterIdentifier,
					"instance_arn":        *instance.DBInstanceArn,
					"instance_identifier": *instance.DBInstanceIdentifier,
				},
			}

			awsRdsInstances = append(awsRdsInstances, instanceTarget)
		}
	}

	return awsRdsInstances, nil
}
