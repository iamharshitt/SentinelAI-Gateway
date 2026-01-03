package policy

import (
	"os"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "policy-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	path := f.Name()
	if _, err := f.WriteString(content); err != nil {
		f.Close()
		os.Remove(path)
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return path
}

func TestLoadPolicy_Valid(t *testing.T) {
	content := `{
  "version": "1",
  "default_action": "allow",
  "rules": [
    {
      "id": "r1",
      "description": "Block secrets",
      "match": { "type": "regex", "patterns": ["secret"] },
      "action": "block",
      "severity": "high"
    }
  ]
}`
	path := writeTempFile(t, content)
	defer os.Remove(path)

	p, err := LoadPolicy(path)
	if err != nil {
		t.Fatalf("LoadPolicy returned error: %v", err)
	}
	if p.Version != "1" {
		t.Fatalf("expected version 1, got %s", p.Version)
	}
	if p.DefaultAction != "allow" {
		t.Fatalf("expected default_action allow, got %s", p.DefaultAction)
	}
	if len(p.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(p.Rules))
	}
}

func TestLoadPolicy_InvalidJSON(t *testing.T) {
	content := `{
  "version": "1",
  "default_action": "allow",
  "rules": [
    { invalid json }
  ]
}`
	path := writeTempFile(t, content)
	defer os.Remove(path)

	_, err := LoadPolicy(path)
	if err == nil {
		t.Fatalf("expected error for invalid json, got nil")
	}
}

func TestLoadPolicy_MissingFields(t *testing.T) {
	content := `{
  "version": "2"
}`
	path := writeTempFile(t, content)
	defer os.Remove(path)

	p, err := LoadPolicy(path)
	if err != nil {
		t.Fatalf("LoadPolicy returned error: %v", err)
	}
	// Missing fields should result in zero-values, not panics
	if p.Version != "2" {
		t.Fatalf("expected version 2, got %s", p.Version)
	}
	if p.DefaultAction != "" {
		t.Fatalf("expected empty default_action, got %s", p.DefaultAction)
	}
}

func TestLoadPolicy_InvalidRegexFile(t *testing.T) {
	content := `{
	"version": "1",
	"default_action": "allow",
	"rules": [
		{"id":"r1","description":"bad regex","match":{"type":"regex","patterns":["(unclosed"] ,"action":"warn","severity":"low"}
	]
}`
	path := writeTempFile(t, content)
	defer os.Remove(path)

	_, err := LoadPolicy(path)
	if err == nil {
		t.Fatalf("expected error for invalid regex, got nil")
	}
}
