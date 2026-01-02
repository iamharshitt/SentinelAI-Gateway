package analyzer

import (
	"regexp"

	"sentinelai/policy"
)

func Analyze(prompt string, p policy.Policy) AnalysisResult {
	result := AnalysisResult{
		Allowed: p.DefaultAction != "block",
		Action:  p.DefaultAction,
	}

	for _, rule := range p.Rules {
		for _, pattern := range rule.Match.Patterns {
			re := regexp.MustCompile(pattern)
			if re.MatchString(prompt) {

				switch rule.Action {

				case "block":
					return AnalysisResult{
						Allowed:  false,
						Action:   "block",
						Reason:   rule.Description,
						Severity: rule.Severity,
					}

				case "warn":
					result.Action = "warn"
					result.Reason = rule.Description
					result.Severity = rule.Severity

				case "redact":
					prompt = re.ReplaceAllString(
						prompt,
						rule.RedactionToken,
					)
					result.Action = "redact"
					result.ModifiedPrompt = prompt
					result.Severity = rule.Severity
				}
			}
		}
	}

	result.Allowed = true
	return result
}
