package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"os"

	"sentinelai/analyzer"
	"sentinelai/policy"
)

type Request struct {
	Prompt string `json:"prompt"`
}

func main() {
	// 1. Load the security policies
	p, err := policy.LoadPolicy("D:/SentinelAI-Gateway/sentinel_policies.json")
	if err != nil {
		panic(err)
	}

	// 2. Initialize the reader for Stdin
	reader := bufio.NewReader(os.Stdin)

	for {
		// 3. Read the 4-byte message length (Native Messaging Protocol)
		lengthBytes := make([]byte, 4)
		_, err := reader.Read(lengthBytes)
		if err != nil {
			return // Exit if the browser closes the pipe
		}

		length := binary.LittleEndian.Uint32(lengthBytes)

		// 4. Read the actual JSON message
		message := make([]byte, length)
		_, err = reader.Read(message)
		if err != nil {
			return
		}

		var req Request
		json.Unmarshal(message, &req)

		// 5. Run the analyzer
		result := analyzer.Analyze(req.Prompt, p)

		// 6. Send the response back (Length prefix + JSON)
		resp, _ := json.Marshal(result)
		binary.Write(os.Stdout, binary.LittleEndian, uint32(len(resp)))
		os.Stdout.Write(resp)
	}
}
