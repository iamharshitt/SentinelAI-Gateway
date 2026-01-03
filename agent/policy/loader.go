package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

// LoadPolicy loads and validates a policy file.
// The policy path is provided by application configuration (flag/env),
// not user-controlled input.
func LoadPolicy(path string) (Policy, error) {
	var p Policy

	// #nosec G304 -- path is controlled by application configuration, not user input
	data, err := os.ReadFile(path)
	if err != nil {
		return p, err
	}

	if err := json.Unmarshal(data, &p); err != nil {
		return p, err
	}

	// Validate regex patterns to avoid runtime panics in analyzer
	for _, r := range p.Rules {
		if r.Match.Type == "regex" {
			for _, pattern := range r.Match.Patterns {
				if _, err := regexp.Compile(pattern); err != nil {
					id := r.ID
					if id == "" {
						id = "<unknown>"
					}
					return p, fmt.Errorf("invalid regex in rule %s: %w", id, err)
				}
			}
		}
	}

	return p, nil
}
