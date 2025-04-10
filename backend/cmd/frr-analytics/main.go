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

	config := configs.LoadConfig()

	sockServer := socket.NewSocket(config["socket"]["UnixSocketLocation"])
	// sockServer := socket.NewSocket("config")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := sockServer.Start(); err != nil {
			fmt.Printf("Error starting socket server: %s\n", err)
			os.Exit(1)
		}
	}()

	// stopAnalyzer := make(chan bool)
	// go runAnalyzerProcess(config, stopAnalyzer)

	<-sigChan
	fmt.Println("\nShutting down...")
}
