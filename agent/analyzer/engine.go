package analyzer

import (
	"regexp"

	"github.com/iamharshitt/SentinelAI-Gateway/policy"
)

type AnalysisResult struct {
	Allowed        bool   `json:"allowed"`
	Action         string `json:"action"`
	ModifiedPrompt string `json:"modified_prompt,omitempty"`
	Reason         string `json:"reason,omitempty"`
	Severity       string `json:"severity,omitempty"`
	RuleID         string `json:"rule_id,omitempty"`
}

func Analyze(prompt string, p policy.Policy) AnalysisResult {
	result := AnalysisResult{
		Allowed: p.DefaultAction != "block",
		Action:  p.DefaultAction,
	}

	for _, rule := range p.Rules {
		// Prefer precompiled patterns when available
		if len(rule.CompiledPatterns) > 0 {
			for _, re := range rule.CompiledPatterns {
				if re.MatchString(prompt) {
					switch rule.Action {
					case "block":
						return AnalysisResult{
							Allowed:  false,
							Action:   "block",
							Reason:   rule.Description,
							Severity: rule.Severity,
							RuleID:   rule.ID,
						}
					case "warn":
						result.Action = "warn"
						result.Reason = rule.Description
						result.Severity = rule.Severity
						result.RuleID = rule.ID
					case "redact":
						prompt = re.ReplaceAllString(prompt, rule.RedactionToken)
						result.Action = "redact"
						result.ModifiedPrompt = prompt
						result.Severity = rule.Severity
						result.RuleID = rule.ID
					}
				}
			}
			continue
		}

		// Fallback: compile on-the-fly (preserve previous test behavior of panicking on invalid regex)
		for _, pattern := range rule.Match.Patterns {
			re, err := regexp.Compile(pattern)
			if err != nil {
				panic(err)
			}
			if re.MatchString(prompt) {
				switch rule.Action {
				case "block":
					return AnalysisResult{
						Allowed:  false,
						Action:   "block",
						Reason:   rule.Description,
						Severity: rule.Severity,
						RuleID:   rule.ID,
					}
				case "warn":
					result.Action = "warn"
					result.Reason = rule.Description
					result.Severity = rule.Severity
					result.RuleID = rule.ID
				case "redact":
					prompt = re.ReplaceAllString(prompt, rule.RedactionToken)
					result.Action = "redact"
					result.ModifiedPrompt = prompt
					result.Severity = rule.Severity
					result.RuleID = rule.ID
				}
			}
		}
	}

	// If we didn't return early from a "block", it's allowed
	result.Allowed = true
	return result
}
