package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {
	fmt.Println("Starting SentinelAI agent (no systray)...")
	path := "./agent/sentinelai.exe"
	cmd := exec.Command(path)
	if err := cmd.Start(); err != nil {
		fmt.Println("Failed to start agent:", err)
		return
	}
	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Println("Agent process exited with error:", err)
		} else {
			fmt.Println("Agent process exited")
		}
	}()

	// Keep launcher alive
	for {
		time.Sleep(1 * time.Hour)
	}
}
