# Explorex Restaurant Platform Golden Path

Independent proof-of-concept inspired by Explorex's public product footprint and DevOps role.

Explorex operates a restaurant technology stack spanning POS/restaurant operations, ordering,
payments, consumer experiences, CRM/reservations, and data. This repository explores the platform
layer needed to ship that kind of system safely and repeatedly.

It does **not** represent Explorex's internal architecture.

## What this demonstrates

A small platform contract and reference AWS architecture for:

- Spring services on EKS
- Angular/static assets through S3 + CloudFront
- ALB ingress
- safe canary/staged rollouts
- automated rollback
- CPU/memory requests and limits
- HPA boundaries
- readiness/liveness probes
- PodDisruptionBudgets
- workload identity / least-privilege IAM
- managed secrets
- metrics, logs, and traces
- latency/availability SLOs
- dependency declarations for queues, caches, and databases
- deployment evidence for incident and release review

## Why this maps to Explorex

The public DevOps role calls for AWS, Linux, Git, Docker, Kubernetes, CodeBuild/CodePipeline,
CloudFront, EC2, ELB, S3, Lambda, CI/CD, Terraform, monitoring, alerting, and tracing.

This repo turns those requirements into an opinionated deployment golden path rather than a generic IaC sample.

## Architecture

```text
Angular/static assets -> S3 + CloudFront
                           |
Users / restaurant clients -> ALB -> EKS application plane
                                      |-- Spring APIs
                                      |-- workers
                                      `-- internal services
                                              |
                                    queue/cache/database

OTel / metrics / logs / alerts across all services
```

## Platform workflow

```text
PR
 |-- unit/static checks
 |-- container build + scan
 |-- platform contract gate
 |-- Terraform/Helm validation
 v
staging -> smoke checks -> production canary -> SLO promotion gate -> promote / rollback
```

## Run locally

```bash
go test ./...
go vet ./...
go run ./cmd/platformctl -contract examples/ordering-api-prod.json
```

## Production-readiness policy

Production workloads must declare safe replica counts, resources, autoscaling, probes, disruption budgets,
progressive rollout, rollback, SLOs, metrics/logs/traces/alerts, managed secrets, workload identity,
least-privilege IAM, restricted egress, dependency readiness, and rollback-compatible DB migrations.

## Restaurant-platform considerations

Restaurant software sees sharp lunch/dinner peaks, and POS/order/payment paths cannot tolerate casual deploy breakage.
This design keeps idempotency, backward-compatible schema changes, safe rollout, fast rollback, and cohort-aware
observability visible in the deployment contract instead of relying on tribal knowledge.

## AWS implementation path

A production version could use EKS, ECR, ALB, S3 + CloudFront, Route 53 + ACM, SQS, ElastiCache/Redis,
Secrets Manager, EKS Pod Identity/IRSA, CloudWatch + OpenTelemetry, Terraform, Helm, CodeBuild/CodePipeline,
and Lambda where operational glue is justified.

## Repository layout

```text
cmd/platformctl/
internal/policy/
examples/
terraform/modules/eks-platform/
terraform/modules/web-edge/
helm/spring-service/
docs/
.github/workflows/
```

## Disclaimer

Independent engineering prototype based only on public Explorex material and the public role description.
