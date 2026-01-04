package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Log writes a compact JSON log line to stderr with timestamp and level.
func Log(level, msg string, fields map[string]interface{}) {
	entry := make(map[string]interface{}, len(fields)+3)
	entry["ts"] = time.Now().UTC().Format(time.RFC3339)
	entry["level"] = level
	entry["msg"] = msg
	for k, v := range fields {
		entry[k] = v
	}
	b, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintln(os.Stderr, "log marshal error:", err)
		return
	}
	fmt.Fprintln(os.Stderr, string(b))
}
