package analyzer

type AnalysisResult struct {
	Allowed        bool   `json:"allowed"`
	Action         string `json:"action"`
	Reason         string `json:"reason,omitempty"`
	Severity       string `json:"severity,omitempty"`
	ModifiedPrompt string `json:"modified_prompt,omitempty"`
}
