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

	// 1. Load the security policies
	p, err := policy.LoadPolicy(policyPath)
	if err != nil {
		log.Fatalf("failed to load policy from %s: %v", policyPath, err)
	}

	// 2. Initialize the reader for Stdin
	reader := bufio.NewReader(os.Stdin)

	for {
		// 3. Read the 4-byte message length (Native Messaging Protocol)
		lengthBytes := make([]byte, 4)
		if _, err := io.ReadFull(reader, lengthBytes); err != nil {
			if err == io.EOF {
				return // graceful exit
			}
			log.Printf("error reading length prefix: %v", err)
			return
		}

		length := binary.LittleEndian.Uint32(lengthBytes)

		// 4. Read the actual JSON message
		if length == 0 {
			log.Printf("received zero-length message, skipping")
			continue
		}

		message := make([]byte, length)
		if _, err := io.ReadFull(reader, message); err != nil {
			log.Printf("error reading message body: %v", err)
			return
		}

		var req Request
		if err := json.Unmarshal(message, &req); err != nil {
			log.Printf("invalid json message: %v", err)
			// respond with an error result
			respObj := map[string]interface{}{"allowed": false, "action": "error", "reason": "invalid json"}
			sendResponse(respObj)
			continue
		}

		// 5. Run the analyzer
		result := analyzer.Analyze(req.Prompt, p)

		// 6. Send the response back (Length prefix + JSON)
		if err := sendResponse(result); err != nil {
			log.Printf("failed to send response: %v", err)
			return
		}
	}
}

func sendResponse(resp interface{}) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	if err := binary.Write(os.Stdout, binary.LittleEndian, uint32(len(data))); err != nil {
		return err
	}
	_, err = os.Stdout.Write(data)
	return err
}
