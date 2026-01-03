package analyzer

import (
	"testing"

	"sentinelai/policy"
)

func TestAnalyze_BlockWarnRedact(t *testing.T) {
	pol := policy.Policy{
		Version:       "1",
		DefaultAction: "allow",
		Rules: []policy.Rule{
			{
				ID:          "r1",
				Description: "Block prompts containing sensitive words",
				Match: policy.Match{
					Type:     "regex",
					Patterns: []string{"(?i)password|secret"},
				},
				Action:   "block",
				Severity: "high",
			},
			{
				ID:          "r2",
				Description: "Warn if prompt contains email addresses",
				Match: policy.Match{
					Type:     "regex",
					Patterns: []string{"[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}"},
				},
				Action:   "warn",
				Severity: "medium",
			},
			{
				ID:          "r3",
				Description: "Redact URLs",
				Match: policy.Match{
					Type:     "regex",
					Patterns: []string{"https?://[^\\s]+"},
				},
				Action:         "redact",
				Severity:       "medium",
				RedactionToken: "[REDACTED_URL]",
			},
		},
	}

	t.Run("allow simple prompt", func(t *testing.T) {
		res := Analyze("hello world", pol)
		if !res.Allowed || res.Action != "allow" {
			t.Fatalf("expected allow, got %+v", res)
		}
	})

	t.Run("block on password", func(t *testing.T) {
		res := Analyze("my password is 1234", pol)
		if res.Allowed || res.Action != "block" {
			t.Fatalf("expected block, got %+v", res)
		}
	})

	t.Run("warn on email", func(t *testing.T) {
		res := Analyze("contact me at test@example.com", pol)
		if !res.Allowed || res.Action != "warn" {
			t.Fatalf("expected warn, got %+v", res)
		}
	})

	t.Run("redact url", func(t *testing.T) {
		res := Analyze("visit https://google.com", pol)
		if !res.Allowed || res.Action != "redact" || res.ModifiedPrompt == "" {
			t.Fatalf("expected redact with modified prompt, got %+v", res)
		}
	})
}
