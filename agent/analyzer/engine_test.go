package analyzer

import (
	"regexp"
	"testing"

	"github.com/iamharshitt/SentinelAI-Gateway/policy"
)

func TestAnalyze_BlockWarnRedact(t *testing.T) {
	pol := policy.Policy{
		Version:       "1",
		DefaultAction: "allow",
		Rules: []policy.Rule{
			{
				ID:          "r_block",
				Description: "block secret",
				Action:      "block",
				Severity:    "high",
				Match:       policy.Match{Type: "regex", Patterns: []string{"secret"}},
			},
			{
				ID:          "r_warn",
				Description: "warn email",
				Action:      "warn",
				Severity:    "medium",
				Match:       policy.Match{Type: "regex", Patterns: []string{"@example\\.com"}},
			},
			{
				ID:             "r_redact",
				Description:    "redact url",
				Action:         "redact",
				Severity:       "low",
				RedactionToken: "[REDACTED]",
				Match:          policy.Match{Type: "regex", Patterns: []string{"https?://\\S+"}},
			},
		},
	}

	// Precompile patterns to exercise the precompiled branch
	for i := range pol.Rules {
		for _, p := range pol.Rules[i].Match.Patterns {
			re := regexp.MustCompile(p)
			pol.Rules[i].CompiledPatterns = append(pol.Rules[i].CompiledPatterns, re)
		}
	}

	// Block case
	res := Analyze("this contains secret data", pol)
	if res.Action != "block" || res.Allowed {
		t.Fatalf("expected block, got %+v", res)
	}

	// Warn case
	res = Analyze("contact me at user@example.com", pol)
	if res.Action != "warn" || !res.Allowed {
		t.Fatalf("expected warn allowed, got %+v", res)
	}

	// Redact case
	res = Analyze("visit https://example.com/page", pol)
	if res.Action != "redact" || res.ModifiedPrompt == "" {
		t.Fatalf("expected redact with modified prompt, got %+v", res)
	}
}
