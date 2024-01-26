# 🔎 Prometheus Service Discovery for AWS RDS

Discover and monitor your AWS RDS clusters effortlessly with this Golang Prometheus HTTP ServiceDiscovery tool. 🚀  
This service solves issues in case you have RDS auto scaling and you need to measure metrics per instance (user cpu time, number of connections etc).

![GitHub](https://img.shields.io/github/license/contember/prometheus-rds-service-discovery)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/contember/prometheus-rds-service-discovery)

## Table of Contents

1. [Setup](#setup)
2. [HTTP Endpoints](#http-endpoints)
3. [IAM Permissions](#iam-permissions)
4. [Example Usage](#example-usage)

---

## Setup

Before you start monitoring your AWS RDS clusters, follow these steps to set up the service discovery:

1. Define the port on which you want to run your app by setting the `$PORT` environment variable. The default fallback port is `8080`.
2. Define the `$AWS_REGION` environment variable with your desired AWS region. Ensure AWS credentials are properly configured if needed. (This service is ECS and EKS ready, so container credentials can also be used.)

To configure clusters for discovery, follow this naming convention:
- `SCRAPER_0_CLUSTER=my-cluster-1`
- `SCRAPER_0_REGION=us-east-1`

You can configure N scraper clusters by incrementing the number in the middle of the environment variable name (e.g., `SCRAPER_1_CLUSTER=my-cluster-2`, `SCRAPER_1_REGION=us-west-2`, and so on).
This service supports multi-region clusters discovery.

---

## HTTP Endpoints

This service provides the following HTTP endpoints:

- `/.health`: Use this endpoint for health checks.
- `/discovery`: This endpoint provides Prometheus HTTP-SD compatible output in JSON format.

---

## IAM Permissions

Ensure your service has the necessary IAM permissions to function correctly.  
You should grant the `DescribeDBInstances` permission in the specific AWS region where your RDS clusters are located.

```terraform

data "aws_iam_policy_document" "rds_service_discovery" {
  statement {
    actions = [
      "rds:DescribeDBInstances",
    ]

    resources = [
      "*"
    ]
  }
}
```

---

## Example Usage 

### Usage with Prometheus PG-Exporter

```yaml
- job_name: 'rds'
  metrics_path: /probe
  params:
    auth_module: [internal]
  http_sd_configs:
  - url: http://rds-service-discovery:8080/discovery
  relabel_configs:
    - source_labels: [__address__]
      target_label: __param_target
    - source_labels: [__param_target]
      target_label: instance
    - target_label: __address__
      replacement: 'pg_exporter.svc.cluster.local'
```
