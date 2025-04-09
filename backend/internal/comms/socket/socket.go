package socket

import (
	"fmt"
	"net"
	"os"
	"sync"
)

var execMutex sync.Mutex

var isRunning bool = true

type Socket struct {
	socketPath string
	listener   net.Listener
	mutex      sync.Mutex
}

func NewSocket(socketPath string) *Socket {
	return &Socket{
		socketPath: socketPath,
		mutex:      sync.Mutex{},
	}
}

func (s *Socket) Start() error {
	os.Remove(s.socketPath)

	l, err := net.ListenUnix("unix", &net.UnixAddr{s.socketPath, "unix"})
	if err != nil {
		return fmt.Errorf("error listening on socket: %w", err)
	}

	s.listener = l

	socketListener = s.listener

	fmt.Printf("Listening on %s ...\n", s.socketPath)

	for isRunning {
		conn, err := l.Accept()
		if err != nil {
			if !isRunning {
				fmt.Println("Socket server shutting down...")
				break
			}
			fmt.Printf("Error accepting connection: %s\n", err.Error())
			continue
		}

		fmt.Println("New client connected")
		s.handleConnection(conn)
	}

	return nil
}

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

	execMutex.Lock()
	defer execMutex.Unlock()

	response := processCommand(message)

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Printf("Error sending response: %s\n", err.Error())
		return
	}
}

func (s *Socket) Close() {
	if s.listener != nil {
		s.listener.Close()
		os.Remove(s.socketPath)
	}
}

var socketListener net.Listener

func exitSocketServer() {
	fmt.Println("Shutting down socket server...")
	isRunning = false

	if socketListener != nil {
		socketListener.Close()
	}

	fmt.Println("Socket server shut down completed")
	os.Exit(0)
}
