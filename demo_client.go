//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"time"
)

func main() {
	cmd := exec.Command(".\\agent\\sentinelai.exe")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("stdin pipe error:", err)
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("stdout pipe error:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("failed to start agent:", err)
		return
	}

	go func() {
		prompts := []string{
			"hello demo",
			"my password is 1234",
			"contact me at test@example.com",
			"visit https://google.com",
		}
		for _, p := range prompts {
			req := map[string]string{"prompt": p}
			data, _ := json.Marshal(req)
			var buf bytes.Buffer
			_ = binary.Write(&buf, binary.LittleEndian, uint32(len(data)))
			buf.Write(data)
			_, err := stdin.Write(buf.Bytes())
			if err != nil {
				fmt.Println("write error:", err)
				return
			}
			time.Sleep(1 * time.Second)
		}
		stdin.Close()
	}()

	for {
		lenBytes := make([]byte, 4)
		if _, err := io.ReadFull(stdout, lenBytes); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read length error:", err)
			break
		}
		length := binary.LittleEndian.Uint32(lenBytes)
		msg := make([]byte, length)
		if _, err := io.ReadFull(stdout, msg); err != nil {
			fmt.Println("read message error:", err)
			break
		}
		fmt.Println("Response:", string(msg))
	}

	_ = cmd.Wait()
}
