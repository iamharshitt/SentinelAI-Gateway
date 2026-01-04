package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"flag"
	"io"
	"os"
	"time"
    "github.com/iamharshitt/SentinelAI-Gateway/policy"
    "github.com/iamharshitt/SentinelAI-Gateway/analyzer"

type Request struct {
	Prompt string `json:"prompt"`
}

func main() {
	policyFlag := flag.String("policy", "", "path to sentinel_policies.json")
	flag.Parse()

	policyPath := *policyFlag
	if policyPath == "" {
		policyPath = os.Getenv("SENTINEL_POLICY")
	}
	if policyPath == "" {
		policyPath = "sentinel_policies.json"
	}

	Serve(os.Stdin, os.Stdout, policyPath)
}

func Serve(reader io.Reader, writer io.Writer, policyPath string) {
	p, err := policy.LoadPolicy(policyPath)
	if err != nil {
		Log("error", "policy_load_failed", map[string]interface{}{"path": policyPath, "error": err.Error()})
		resp := map[string]interface{}{
			"allowed": false,
			"action":  "error",
			"reason":  "policy_load_failed",
		}
		_ = sendResponseTo(writer, resp)
		return
	}
	Log("info", "policy_loaded", map[string]interface{}{"path": policyPath, "rules": len(p.Rules)})

	// Start metrics server (Prometheus)
	if err := ServeMetrics(":9090"); err == nil {
		Log("info", "metrics_listening", map[string]interface{}{"addr": ":9090"})
	}

	bufReader := bufio.NewReader(reader)
	lengthBytes := make([]byte, 4) // Moved outside loop for efficiency

	for {
		if _, err := io.ReadFull(bufReader, lengthBytes); err != nil {
			return
		}

		length := binary.LittleEndian.Uint32(lengthBytes)
		if length == 0 {
			continue
		}

		message := make([]byte, length)
		if _, err := io.ReadFull(bufReader, message); err != nil {
			return
		}

		var req Request
		if err := json.Unmarshal(message, &req); err != nil {
			Log("warn", "invalid_json", map[string]interface{}{"error": err.Error()})
			resp := map[string]interface{}{
				"allowed": false,
				"action":  "error",
				"reason":  "invalid_json",
			}
			_ = sendResponseTo(writer, resp)
			continue
		}

		// Instrumentation: count request and measure latency
		requestsTotal.Inc()
		start := time.Now()
		result := analyzer.Analyze(req.Prompt, p)
		dur := time.Since(start).Seconds()
		requestLatency.Observe(dur)
		actionCounter.WithLabelValues(result.Action).Inc()

		// Log the result (do not include raw prompt to avoid leaking sensitive information)
		Log("info", "analyze_result", map[string]interface{}{"action": result.Action, "allowed": result.Allowed, "rule_id": result.RuleID, "severity": result.Severity, "duration_s": dur})

		if err := sendResponseTo(writer, result); err != nil {
			Log("error", "send_response_failed", map[string]interface{}{"error": err.Error()})
			return
		}
	}
}

func sendResponseTo(writer io.Writer, resp interface{}) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	dataLen := len(data)
	if err := binary.Write(writer, binary.LittleEndian, uint32(dataLen)); err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}
