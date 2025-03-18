package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/ba2025-ysmprc/frr-tui/internal/ui/views"
	"github.com/ba2025-ysmprc/frr-tui/internal/utils"

	"github.com/ba2025-ysmprc/frr-tui/internal/batfish"
)

func startBatfish() error {
	// Start Batfish using Docker Compose
	cmd := exec.Command("./scripts/start-batfish.sh")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start Batfish: %v", err)
	}
	return nil
}

func stopBatfish() error {
	// Stop Batfish using Docker Compose
	cmd := exec.Command("./scripts/stop-batfish.sh")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop Batfish: %v", err)
	}
	return nil
}

func main() {
	// Start Batfish
	if err := startBatfish(); err != nil {
		log.Fatalf("Failed to start Batfish: %v", err)
	}

	// Ensure Batfish is stopped when the TUI exits
	defer func() {
		if err := stopBatfish(); err != nil {
			log.Printf("Failed to stop Batfish: %v", err)
		} else {
			log.Println("Batfish stopped successfully.")
		}
	}()

	// Catch termination signals (e.g., Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down...")
		if err := stopBatfish(); err != nil {
			log.Printf("Failed to stop Batfish: %v", err)
		}
		os.Exit(0)
	}()

	// Initialize Batfish client
	bfClient := batfish.NewBatfishClient()

	// Simulate network monitoring
	for {
		anomalyDetected, err := utils.MonitorNetwork()
		if err != nil {
			log.Fatalf("Failed to monitor network: %v", err)
		}

		if anomalyDetected {
			fmt.Println("Anomaly detected! Analyzing with Batfish...")

			// Upload snapshot (replace with your snapshot path)
			snapshotPath := "configs/snapshots"
			snapshotName := "snapshot"
			if err := bfClient.UploadSnapshot(snapshotPath, snapshotName); err != nil {
				log.Fatalf("Failed to upload snapshot: %v", err)
			}

			// Run analysis (replace with your question)
			result, err := bfClient.RunAnalysis(snapshotName, "bgpSessionStatus")
			if err != nil {
				log.Fatalf("Failed to run analysis: %v", err)
			}

			// Display results in TUI
			ui.DisplayResults(result)
		}
	}
}
