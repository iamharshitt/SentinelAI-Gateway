package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("SentinelAI")
	systray.SetTooltip("Secure AI Gateway Running")

	go func() {
		fmt.Println("Starting SentinelAI agent...")
		path := "./agent/sentinelai.exe"
		cmd := exec.Command(path)
		err := cmd.Start()
		if err != nil {
			fmt.Println("Failed to start agent:", err)
			return
		}
		if err := cmd.Wait(); err != nil {
			fmt.Println("Agent process exited with error:", err)
		}
	}()

	// Keep tray alive
	go func() {
		for {
			time.Sleep(1 * time.Hour)
		}
	}()
}

func onExit() {
	fmt.Println("SentinelAI Tray Exited")
}
