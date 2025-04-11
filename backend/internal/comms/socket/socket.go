package socket

import (
	"fmt"
	"net"
	"os"
	"sync"
)

// Global mutex to ensure synchronous execution
var execMutex sync.Mutex

// Global variable to track running state
var isRunning bool = true

// Socket represents a Unix socket server
type Socket struct {
	socketPath string
	listener   net.Listener
	mutex      sync.Mutex
}

// NewSocket creates a new Socket instance
func NewSocket(socketPath string) *Socket {
	return &Socket{
		socketPath: socketPath,
		mutex:      sync.Mutex{},
	}
}

// Start begins listening on the Unix socket
func (s *Socket) Start() error {
	// Clean up any existing socket file
	os.Remove(s.socketPath)

	// Create and listen on the Unix socket
	l, err := net.ListenUnix("unix", &net.UnixAddr{s.socketPath, "unix"})
	if err != nil {
		return fmt.Errorf("error listening on socket: %w", err)
	}

	s.listener = l

	// Set global reference to listener for exit command
	socketListener = s.listener

	fmt.Printf("Listening on %s ...\n", s.socketPath)

	// Accept connections while isRunning is true
	for isRunning {
		conn, err := l.Accept()
		if err != nil {
			// Check if we're shutting down
			if !isRunning {
				fmt.Println("Socket server shutting down...")
				break
			}
			fmt.Printf("Error accepting connection: %s\n", err.Error())
			continue
		}

		fmt.Println("New client connected")
		// Handle connection
		s.handleConnection(conn)
	}

	return nil
}

// handleConnection processes client connections
func (s *Socket) handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Error reading from connection: %s\n", err.Error())
		return
	}

	message := string(buf[0:n])
	fmt.Println("Received message:", message)

	// Lock mutex to ensure synchronous processing
	execMutex.Lock()
	defer execMutex.Unlock()

	// Process the command based on the message
	response := processCommand(message)

	// Send the response back
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Printf("Error sending response: %s\n", err.Error())
		return
	}
}

// Close shuts down the socket
func (s *Socket) Close() {
	if s.listener != nil {
		s.listener.Close()
		os.Remove(s.socketPath)
	}
}

// Global reference to the socket listener for the exit command
var socketListener net.Listener

// exitSocketServer gracefully shuts down the socket server
func exitSocketServer() {
	fmt.Println("Shutting down socket server...")
	isRunning = false

	// Close the listener to stop accepting new connections
	if socketListener != nil {
		socketListener.Close()
	}

	// Note: We don't need to call os.Exit(0) here because we want the main function
	// to handle the clean shutdown. The main function already calls sockServer.Close()
	// when it receives a signal to shut down.

	fmt.Println("Socket server shut down completed")
	os.Exit(0)
}
