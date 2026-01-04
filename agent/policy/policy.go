package policy

import "regexp"

type Match struct {
	Type     string   `json:"type"`
	Patterns []string `json:"patterns"`
}

type Rule struct {
	ID             string `json:"id"`
	Description    string `json:"description"`
	Match          Match  `json:"match"`
	Action         string `json:"action"`
	Severity       string `json:"severity"`
	RedactionToken string `json:"redaction_token,omitempty"`

	// Add this field!
	// We use 'json:"-"' so it doesn't try to look for it in the JSON file
	CompiledPatterns []*regexp.Regexp `json:"-"`
}

type Policy struct {
	Version       string `json:"version"`
	DefaultAction string `json:"default_action"`
	Rules         []Rule `json:"rules"`
}
