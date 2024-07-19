package service_scraper

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go/aws/arn"
	log "github.com/sirupsen/logrus"
)

type ScrapingTarget struct {
	ClusterArn           string
	ClusterAssumeRoleArn *string
	ClusterNumber        int
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
	instances := make([]ScraperDiscoveredTarget, 0)

	// iterate over all registered clusters
	for _, cluster := range s.Targets {
		awsRdsClient, err := s.getClusterAwsClient(ctx, cluster)
		if err != nil {
			log.Errorf("Error creating AWS client for cluster %s %v", cluster.ClusterArn, err)
			continue
		}

		// describe all instances in the cluster with expected cluster ARN
		resp, err := awsRdsClient.DescribeDBInstances(
			ctx,
			&rds.DescribeDBInstancesInput{
				Filters: []types.Filter{
					{Name: aws.String("db-cluster-id"), Values: []string{cluster.ClusterArn}},
				},
			},
		)
		if err != nil {
			log.Errorf("Error describing RDS cluster %s %v", cluster.ClusterArn, err)
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

			instances = append(instances, instanceTarget)
		}
	}

	return instances, nil
}

func (s *Scraper) getClusterAwsClient(ctx context.Context, target ScrapingTarget) (*rds.Client, error) {
	conf, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	// basic AWS client without assuming role
	if target.ClusterAssumeRoleArn == nil {
		return rds.NewFromConfig(conf), nil
	}

	targetArn, err := arn.Parse(target.ClusterArn)
	if err != nil {
		return nil, fmt.Errorf("invalid RDS cluster ARN: %s", target.ClusterArn)
	}

	// create a new AWS client with assumed role
	stsClient := sts.NewFromConfig(conf)
	assumedRoleCred := stscreds.NewAssumeRoleProvider(stsClient, *target.ClusterAssumeRoleArn)

	assumedRoleConfig := conf
	assumedRoleConfig.Credentials = aws.NewCredentialsCache(assumedRoleCred)

	// this must be here because in AWS RDS API there is error,
	// you must explicitly set region where cluster is located when assuming role,
	// or you will receive weird "The parameter Filter: db-cluster-id is not a valid identifier." error message
	// idk why, but it is what it is
	assumedRoleConfig.Region = targetArn.Region

	return rds.NewFromConfig(assumedRoleConfig), nil
}
