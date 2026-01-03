package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"sentinelai/analyzer"
	"sentinelai/policy"
)

type Request struct {
	Prompt string `json:"prompt"`
}

func main() {
	// Accept policy path via flag or environment variable
	policyFlag := flag.String("policy", "", "path to sentinel_policies.json")
	flag.Parse()

	policyPath := *policyFlag
	if policyPath == "" {
		policyPath = os.Getenv("SENTINEL_POLICY")
	}
	if policyPath == "" {
		policyPath = "D:/SentinelAI-Gateway/sentinel_policies.json"
	}

	// Serve on the process stdio
	Serve(os.Stdin, os.Stdout, policyPath)
}

// Serve runs the native messaging loop reading length-prefixed JSON
// messages from reader and writing length-prefixed JSON responses to writer.
func Serve(reader io.Reader, writer io.Writer, policyPath string) {
	// 1. Load the security policies
	p, err := policy.LoadPolicy(policyPath)
	if err != nil {
		log.Fatalf("failed to load policy from %s: %v", policyPath, err)
	}

	bufReader := bufio.NewReader(reader)

	for {
		lengthBytes := make([]byte, 4)
		if _, err := io.ReadFull(bufReader, lengthBytes); err != nil {
			if err == io.EOF {
				return // graceful exit
			}
			log.Printf("error reading length prefix: %v", err)
			return
		}

		length := binary.LittleEndian.Uint32(lengthBytes)

		if length == 0 {
			log.Printf("received zero-length message, skipping")
			continue
		}

		message := make([]byte, length)
		if _, err := io.ReadFull(bufReader, message); err != nil {
			log.Printf("error reading message body: %v", err)
			return
		}

		var req Request
		if err := json.Unmarshal(message, &req); err != nil {
			log.Printf("invalid json message: %v", err)
			respObj := map[string]interface{}{"allowed": false, "action": "error", "reason": "invalid json"}
			if err := sendResponseTo(writer, respObj); err != nil {
				log.Printf("failed to send response: %v", err)
				return
			}
			continue
		}

		// 5. Run the analyzer
		result := analyzer.Analyze(req.Prompt, p)

		// 6. Send the response back (Length prefix + JSON)
		if err := sendResponseTo(writer, result); err != nil {
			log.Printf("failed to send response: %v", err)
			return
		}
	}
}

func sendResponseTo(writer io.Writer, resp interface{}) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	if err := binary.Write(writer, binary.LittleEndian, uint32(len(data))); err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}
