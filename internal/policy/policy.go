package policy

import "fmt"

type SLO struct {
	AvailabilityPercent float64 `json:"availability_percent"`
	P95LatencyMS        int     `json:"p95_latency_ms"`
}

type Observability struct {
	Metrics bool `json:"metrics"`
	Logs    bool `json:"logs"`
	Traces  bool `json:"traces"`
	Alerts  bool `json:"alerts"`
}

type Rollout struct {
	Strategy        string `json:"strategy"`
	RollbackEnabled bool   `json:"rollback_enabled"`
}

type Contract struct {
	Service             string        `json:"service"`
	Environment         string        `json:"environment"`
	Replicas            int           `json:"replicas"`
	CPURequest          string        `json:"cpu_request"`
	CPULimit            string        `json:"cpu_limit"`
	MemoryRequest       string        `json:"memory_request"`
	MemoryLimit         string        `json:"memory_limit"`
	MinReplicas         int           `json:"min_replicas"`
	MaxReplicas         int           `json:"max_replicas"`
	ReadinessProbe      bool          `json:"readiness_probe"`
	LivenessProbe       bool          `json:"liveness_probe"`
	PDB                  bool          `json:"pod_disruption_budget"`
	SecretSource        string        `json:"secret_source"`
	WorkloadIdentity    bool          `json:"workload_identity"`
	IAMWildcards        []string      `json:"iam_wildcards"`
	PublicAdminIngress  bool          `json:"public_admin_ingress"`
	RestrictedEgress    bool          `json:"restricted_egress"`
	DependenciesReady   bool          `json:"dependencies_ready"`
	MigrationStrategy   string        `json:"migration_strategy"`
	IdempotencyRequired bool          `json:"idempotency_required"`
	SLO                 SLO           `json:"slo"`
	Observability       Observability `json:"observability"`
	Rollout             Rollout       `json:"rollout"`
}

type Finding struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

type Result struct {
	Allowed  bool      `json:"allowed"`
	Findings []Finding `json:"findings,omitempty"`
}

func Evaluate(c Contract) Result {
	var f []Finding
	prod := c.Environment == "production"

	if prod {
		if c.Replicas < 2 {
			f = append(f, Finding{"high", "production stateless services require at least 2 replicas"})
		}
		if c.CPURequest == "" || c.CPULimit == "" || c.MemoryRequest == "" || c.MemoryLimit == "" {
			f = append(f, Finding{"high", "CPU and memory requests/limits are required"})
		}
		if c.MinReplicas < 2 || c.MaxReplicas < c.MinReplicas {
			f = append(f, Finding{"high", "autoscaling bounds are invalid"})
		}
		if !c.ReadinessProbe || !c.LivenessProbe {
			f = append(f, Finding{"high", "readiness and liveness probes are required"})
		}
		if !c.PDB {
			f = append(f, Finding{"medium", "PodDisruptionBudget is required"})
		}
		if !c.WorkloadIdentity {
			f = append(f, Finding{"critical", "workload identity is required"})
		}
		if c.SecretSource == "" || c.SecretSource == "plaintext" {
			f = append(f, Finding{"critical", "managed secret source required"})
		}
		if len(c.IAMWildcards) > 0 {
			f = append(f, Finding{"critical", "wildcard IAM permissions are not allowed"})
		}
		if c.PublicAdminIngress {
			f = append(f, Finding{"critical", "administrative ingress cannot be public"})
		}
		if !c.RestrictedEgress {
			f = append(f, Finding{"high", "production egress must be explicitly bounded"})
		}
		if !c.DependenciesReady {
			f = append(f, Finding{"high", "declared dependencies must pass readiness checks"})
		}
		if c.SLO.AvailabilityPercent < 99.9 || c.SLO.P95LatencyMS <= 0 {
			f = append(f, Finding{"high", "availability and latency SLOs are required"})
		}
		if !c.Observability.Metrics || !c.Observability.Logs || !c.Observability.Traces || !c.Observability.Alerts {
			f = append(f, Finding{"high", "metrics, logs, traces, and alerts are required"})
		}
	}

	switch c.Rollout.Strategy {
	case "canary", "blue-green", "staged":
	default:
		f = append(f, Finding{"high", fmt.Sprintf("unsupported rollout strategy %q", c.Rollout.Strategy)})
	}
	if !c.Rollout.RollbackEnabled {
		f = append(f, Finding{"high", "rollback must be enabled"})
	}

	switch c.MigrationStrategy {
	case "expand-contract", "backward-compatible":
	default:
		f = append(f, Finding{"high", "database migration must remain rollback-compatible"})
	}

	allowed := true
	for _, x := range f {
		if x.Severity == "high" || x.Severity == "critical" {
			allowed = false
		}
	}
	return Result{Allowed: allowed, Findings: f}
}
