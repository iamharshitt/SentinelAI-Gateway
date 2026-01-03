package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os/exec"
)

func sendPrompt(prompt string) {
	req := map[string]string{"prompt": prompt}
	data, err := json.Marshal(req)
	if err != nil {
		fmt.Println("json marshal error:", err)
		return
	}

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, uint32(len(data))); err != nil {
		fmt.Println("binary write error:", err)
		return
	}
	buf.Write(data)

	cmd := exec.Command(".\\agent\\sentinelai.exe")
	cmd.Stdin = &buf

	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(out) < 4 {
		fmt.Println("Invalid response from agent")
		return
	}

	length := binary.LittleEndian.Uint32(out[:4])
	if int(4+length) > len(out) {
		fmt.Println("Truncated response from agent")
		return
	}
	resp := out[4 : 4+length]

	fmt.Println("Prompt:", prompt)
	fmt.Println("Response:", string(resp))
	fmt.Println("----")
}

func main() {
	sendPrompt("hello world")
	sendPrompt("my password is 1234")
	sendPrompt("contact me at test@example.com")
	sendPrompt("visit https://google.com")
}
