package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("[entrypoint] ROAR starting...")

	mode := os.Getenv("ROAR_MODE")
	if mode == "" {
		mode = "api"
	}

	fmt.Printf("[entrypoint] Mode: %s\n", mode)

	fmt.Println("[entrypoint] Loading configuration...")
	os.Setenv("ROAR_ROOT", "/app")
	os.Setenv("ROAR_RESULTS_DIR", "/app/new-results")
	os.Setenv("ROAR_CONFIG_FILE", "/app/roar.yaml")
	os.Setenv("ROAR_ENV", "docker")
	fmt.Println("[entrypoint] Configuration loaded successfully")

	fmt.Println("[entrypoint] Database schema initialization delegated to API/Bot startup")

	if os.Getenv("RUN_SETUP") == "true" {
		fmt.Println("[entrypoint] RUN_SETUP=true -> executing check-tools (Go)")
		cmd := exec.Command("/usr/local/bin/roar", "check-tools")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()
		if err := cmd.Run(); err != nil {
			fmt.Println("[entrypoint] check-tools finished with warnings")
		} else {
			fmt.Println("[entrypoint] check-tools completed successfully")
		}
	}

	resultsDir := os.Getenv("ROAR_RESULTS_DIR")
	if resultsDir == "" {
		resultsDir = "/app/new-results"
	}
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		fmt.Printf("[entrypoint] Warning: Failed to create results directory: %v\n", err)
	}

	if _, err := os.Stat("/usr/local/bin/roar"); err != nil {
		fmt.Fprintf(os.Stderr, "[entrypoint] Error: ROAR binary not found at /usr/local/bin/roar\n")
		os.Exit(1)
	}

	var cmd *exec.Cmd
	switch mode {
	case "api", "discord", "both":
		fmt.Println("[entrypoint] Launching API Server...")
		cmd = exec.Command("/usr/local/bin/roar", "api")
	default:
		fmt.Fprintf(os.Stderr, "[entrypoint] Error: Invalid ROAR_MODE '%s'\n", mode)
		fmt.Fprintf(os.Stderr, "[entrypoint] Valid modes: api\n")
		os.Exit(1)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "[entrypoint] Error: %v\n", err)
		os.Exit(1)
	}
}
