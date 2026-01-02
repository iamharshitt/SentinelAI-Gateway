package policy

type Match struct {
	Type     string   `json:"type"`
	Patterns []string `json:"patterns"`
}

type Rule struct {
	ID             string `json:"id"`
	Description    string `json:"description"`
	Match          Match  `json:"match"`
	Action         string `json:"action"` // block | warn | redact
	Severity       string `json:"severity"`
	RedactionToken string `json:"redaction_token,omitempty"`
}

type Policy struct {
	Version       string `json:"version"`
	DefaultAction string `json:"default_action"`
	Rules         []Rule `json:"rules"`
}
