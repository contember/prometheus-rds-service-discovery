package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

func scrapeClusters(config ScrapingConfig) []Target {
	var awsRdsClusters = []Target{}

	for _, target := range config.Targets {
		awsSession, _ := awsConfig.LoadDefaultConfig(
			context.TODO(),
			awsConfig.WithRegion(target.Region),
		)

		awsRds := rds.NewFromConfig(awsSession)

		// filters all instances that contains cluster id that we want
		filters := []types.Filter{
			{
				Name:   aws.String("db-cluster-id"),
				Values: []string{target.Cluster},
			},
		}

		describeInput := &rds.DescribeDBInstancesInput{
			Filters: filters,
		}

		resp, err := awsRds.DescribeDBInstances(context.TODO(), describeInput)

		if err != nil {
			fmt.Printf("Error describing cluster %s: %v\n", target.Cluster, err)
			continue
		}

		for _, cluster := range resp.DBInstances {
			var clusterDsn = *cluster.Endpoint.Address + ":" + fmt.Sprint(*cluster.Endpoint.Port)

			var instanceTarget = Target{
				Targets: []string{clusterDsn},
				Labels: map[string]string{
					"cluster_region":      target.Region,
					"cluster_identifier":  *cluster.DBClusterIdentifier,
					"instance_identifier": *cluster.DBInstanceIdentifier,
				},
			}

			awsRdsClusters = append(awsRdsClusters, instanceTarget)
		}
	}

	return awsRdsClusters
}
