package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"flag"
	"io"
	"os"

	"sentinelai/analyzer"
	"sentinelai/policy"
)

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
		resp := map[string]interface{}{
			"allowed": false,
			"action":  "error",
			"reason":  "policy_load_failed",
		}
		_ = sendResponseTo(writer, resp)
		return
	}

	bufReader := bufio.NewReader(reader)

	for {
		lengthBytes := make([]byte, 4)
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
			resp := map[string]interface{}{
				"allowed": false,
				"action":  "error",
				"reason":  "invalid_json",
			}
			_ = sendResponseTo(writer, resp)
			continue
		}

		result := analyzer.Analyze(req.Prompt, p)
		if err := sendResponseTo(writer, result); err != nil {
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
	if dataLen > int(^uint32(0)) {
		return io.ErrShortBuffer
	}

	if err := binary.Write(writer, binary.LittleEndian, uint32(dataLen)); err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}
