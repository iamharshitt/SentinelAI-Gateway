package analyzer

import (
	"testing"

	"sentinelai/policy"
)

func TestAnalyze_PanicOnInvalidRegex(t *testing.T) {
	pol := policy.Policy{
		Version:       "1",
		DefaultAction: "allow",
		Rules: []policy.Rule{
			{
				ID:          "r1",
				Description: "bad regex",
				Match: policy.Match{
					Type:     "regex",
					Patterns: []string{"(unclosed"},
				},
				Action:   "warn",
				Severity: "low",
			},
		},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic from invalid regex, but no panic occurred")
		}
	}()

	_ = Analyze("test", pol)
}
