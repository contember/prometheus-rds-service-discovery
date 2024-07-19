package service_scraper

import (
	"os"
	"testing"
)

func TestScrapeWithoutAssumeRoles(t *testing.T) {
	arnA := "arn:aws:rds:eu-central-1:123456789012:cluster:first-cluster"
	arnB := "arn:aws:rds:eu-central-1:123456789012:cluster:second-cluster"

	_ = os.Setenv("SCRAPER_0_CLUSTER_ARN", arnA)
	_ = os.Setenv("SCRAPER_1_CLUSTER_ARN", arnB)

	targets, err := scraperTargetsEnvParser()
	if err != nil {
		t.Fatalf("Error parsing targets: %v", err)
	}

	if len(targets) != 2 {
		t.Fatalf("Expected 2 targets, got %d", len(targets))
	}

	if targets[0].ClusterArn != arnA {
		t.Fatalf("Expected first target ARN to be %s, got %s", arnA, targets[0].ClusterArn)
	}

	if targets[1].ClusterArn != arnB {
		t.Fatalf("Expected first target ARN to be %s, got %s", arnB, targets[1].ClusterArn)
	}
}

func TestScrapeWithAssumeRoles(t *testing.T) {
	arn := "arn:aws:rds:eu-central-1:123456789012:cluster:first-cluster"
	arnAssume := "arn:aws:iam::123456789013:role/some-fancy-role"

	_ = os.Setenv("SCRAPER_0_CLUSTER_ARN", arn)
	_ = os.Setenv("SCRAPER_0_ASSUME_ROLE_ARN", arnAssume)

	targets, err := scraperTargetsEnvParser()
	if err != nil {
		t.Fatalf("Error parsing targets: %v", err)
	}

	if len(targets) != 1 {
		t.Fatalf("Expected 1 targets, got %d", len(targets))
	}

	if targets[0].ClusterArn != arn {
		t.Fatalf("Expected first target ARN to be %s, got %s", arn, targets[0].ClusterArn)
	}

	if targets[0].ClusterAssumeRoleArn == nil || *targets[0].ClusterAssumeRoleArn != arnAssume {
		t.Fatalf("Expected first target ARN to be %s", arnAssume)
	}
}
