package socket

import (
	"context"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	SocketPath string
	BufferSize int
}

func DefaultConfig() Config {
	return Config{
		SocketPath: "/tmp/detector.sock",
		BufferSize: 4096,
	}
}

type Server struct {
	config     Config
	listener   net.Listener
	handlers   map[string]Handler
	mu         sync.RWMutex
	wg         sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
}

type Handler func([]byte) ([]byte, error)

func CreateNewSocket(config Config) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		config:     config,
		handlers:   make(map[string]Handler),
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

func (s *Server) getHandler(command string) (Handler, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	handler, ok := s.handlers[command]
	return handler, ok
}

func (s *Server) Start() error {
	dir := filepath.Dir(s.config.SocketPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err := os.RemoveAll(s.config.SocketPath); err != nil {
		return err
	}

	listener, err := net.Listen("unix", s.config.SocketPath)
	if err != nil {
		return err
	}
	s.listener = listener

	if err := os.Chmod(s.config.SocketPath, 0660); err != nil {
		log.Printf("Failed to set socket permissions: %v", err)
	}

	s.wg.Add(1)
	go s.acceptConnections()

	log.Printf("Socket server started on %s", s.config.SocketPath)
	return nil
}

func (s *Server) acceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.ctx.Done():
					return
				default:
					log.Printf("Error accepting connection: %v", err)
					continue
				}
			}

			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	buffer := make([]byte, s.config.BufferSize)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Printf("Error reading from connection: %v", err)
		return
	}

	if n < 4 {
		log.Printf("Message too short")
		return
	}

	command := string(buffer[:4])
	payload := buffer[4:n]
	handler, ok := s.getHandler(command)
	if !ok {
		log.Printf("Unknown command: %s", command)
		return
	}

	response, err := handler(payload)
	if err != nil {
		log.Printf("Error handling command %s: %v", command, err)
		return
	}

	if _, err := conn.Write(response); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (s *Server) Stop() error {
	s.cancelFunc()

	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			return err
		}
	}

	s.wg.Wait()

	if err := os.RemoveAll(s.config.SocketPath); err != nil {
		return err
	}

	log.Printf("Socket server stopped")
	return nil
}
