package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ba2025-ysmprc/frr-tui/backend/configs"
	socket "github.com/ba2025-ysmprc/frr-tui/backend/internal/comms/socket"
)

func main() {

	configs.LoadConfig()

	os.Exit(0)
	sockServer := socket.NewSocket("/tmp/unixsock.sock")

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the socket server in a goroutine
	go func() {
		if err := sockServer.Start(); err != nil {
			fmt.Printf("Error starting socket server: %s\n", err)
			os.Exit(1)
		}
	}()

	fmt.Println("Socket server running. Press Ctrl+C to exit.")

	// Wait for signal
	<-sigChan
	fmt.Println("\nShutting down...")
	sockServer.Close()
	fmt.Println("Server stopped.")
}
