package policy

import "testing"

func safe() Contract {
	return Contract{
		Service: "ordering-api", Environment: "production", Replicas: 3,
		CPURequest: "500m", CPULimit: "2", MemoryRequest: "1Gi", MemoryLimit: "3Gi",
		MinReplicas: 3, MaxReplicas: 20,
		ReadinessProbe: true, LivenessProbe: true, PDB: true,
		SecretSource: "aws-secrets-manager", WorkloadIdentity: true,
		IAMWildcards: nil, PublicAdminIngress: false, RestrictedEgress: true,
		DependenciesReady: true, MigrationStrategy: "expand-contract", IdempotencyRequired: true,
		SLO: SLO{AvailabilityPercent: 99.95, P95LatencyMS: 350},
		Observability: Observability{Metrics: true, Logs: true, Traces: true, Alerts: true},
		Rollout: Rollout{Strategy: "canary", RollbackEnabled: true},
	}
}

func TestSafeContract(t *testing.T) {
	r := Evaluate(safe())
	if !r.Allowed {
		t.Fatalf("expected allowed: %+v", r.Findings)
	}
}

func TestWildcardIAMRejected(t *testing.T) {
	c := safe()
	c.IAMWildcards = []string{"s3:*"}
	if Evaluate(c).Allowed {
		t.Fatal("expected IAM rejection")
	}
}

func TestUnsafeMigrationRejected(t *testing.T) {
	c := safe()
	c.MigrationStrategy = "drop-and-recreate"
	if Evaluate(c).Allowed {
		t.Fatal("expected migration rejection")
	}
}
