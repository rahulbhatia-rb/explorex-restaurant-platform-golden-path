# Architecture Notes

## Workload classes

A restaurant platform should distinguish:
- synchronous ordering/payment APIs
- asynchronous workers
- web/frontend delivery
- reporting/search/data paths
- internal/admin services

Their scaling and SLO profiles are different and should not share one generic policy.

## Peak windows

Restaurant traffic often has predictable lunch/dinner peaks. Capacity planning can use historical peak concurrency, queue depth, p95/p99 latency, CPU/memory saturation, outlet cohort load, and node provisioning time.

## Deployment safety

For order/payment-critical services:
- canary before broad rollout
- health + SLO checks before promotion
- backward-compatible schema changes
- idempotent mutation APIs
- rapid rollback
- release metadata in traces/logs

## Observability

Every request should be traceable across edge -> API -> queue/cache/database -> downstream service. Dashboards should support slicing by environment, service, version, and customer/outlet cohort.
