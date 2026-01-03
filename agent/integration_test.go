package main

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"
)

func writeMessage(w io.Writer, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint32(len(data))); err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func readMessage(r io.Reader) ([]byte, error) {
	var lenb [4]byte
	if _, err := io.ReadFull(r, lenb[:]); err != nil {
		return nil, err
	}
	length := binary.LittleEndian.Uint32(lenb[:])
	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}
	return data, nil
}

func TestServe_NativeMessagingLoop(t *testing.T) {
	// write a minimal policy file
	policyContent := `{
  "version": "1",
  "default_action": "allow",
  "rules": [
    {"id":"r1","description":"block pw","match":{"type":"regex","patterns":["password"]},"action":"block","severity":"high"},
    {"id":"r2","description":"warn email","match":{"type":"regex","patterns":["[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}"]},"action":"warn","severity":"medium"},
    {"id":"r3","description":"redact url","match":{"type":"regex","patterns":["https?://[^\\s]+"]},"action":"redact","severity":"medium","redaction_token":"[REDACTED_URL]"}
  ]
}`

	tmpFile, err := os.CreateTemp("", "policy-*.json")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(policyContent); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	tmpFile.Close()

	pr, pw := io.Pipe()
	rr, rw := io.Pipe()

	// run Serve in goroutine
	go Serve(pr, rw, tmpFile.Name())

	// small delay to ensure goroutine is ready
	time.Sleep(50 * time.Millisecond)

	// client writes to pw, reads from rr
	tests := []struct {
		prompt       string
		expectAction string
	}{
		{"hello world", "allow"},
		{"my password is 1234", "block"},
		{"contact me at test@example.com", "warn"},
		{"visit https://google.com", "redact"},
	}

	for _, tc := range tests {
		if err := writeMessage(pw, map[string]string{"prompt": tc.prompt}); err != nil {
			t.Fatalf("write message: %v", err)
		}
		data, err := readMessage(rr)
		if err != nil {
			t.Fatalf("read message: %v", err)
		}
		var out map[string]interface{}
		if err := json.Unmarshal(data, &out); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		act, _ := out["action"].(string)
		if act == "" && tc.expectAction == "allow" {
			// analyzer returns default "allow" but may set action to "allow"
			// accept both
			continue
		}
		if act != tc.expectAction {
			t.Fatalf("prompt %q expected action %s got %v", tc.prompt, tc.expectAction, act)
		}
	}

	// close pipes
	pw.Close()
	rr.Close()
}
