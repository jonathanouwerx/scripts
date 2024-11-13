package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mycli <script_name>")
		os.Exit(1)
	}

	// TODO: modify this line to be the location of the scripts you want on your own machine
	scriptDir := filepath.Join(os.Getenv("HOME"), "code", "personal", "scripts", "scripts_bin")

	scriptName := os.Args[1]
	scriptPath := filepath.Join(scriptDir, scriptName+".sh")

	// Check if the script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		fmt.Printf("Script %s not found in %s\n", scriptName, scriptDir)
		os.Exit(1)
	}

	// Execute the script
	cmd := exec.Command(scriptPath, os.Args[2:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running script %s: %v\n", scriptName, err)
		os.Exit(1)
	}
}
