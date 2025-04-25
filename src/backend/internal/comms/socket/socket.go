package socket

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/configs"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/analyzer"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/logger"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"google.golang.org/protobuf/proto"
)

var execMutex sync.Mutex

var isRunning bool = true

type Socket struct {
	socketPath string
	listener   net.Listener
	mutex      sync.Mutex
	metrics    *frrProto.FullFRRData
	analyzer   *analyzer.Analyzer
	logger     *logger.Logger
}

func NewSocket(config configs.SocketConfig, metrics *frrProto.FullFRRData, analyzer *analyzer.Analyzer, logger *logger.Logger) *Socket {
	return &Socket{
		socketPath: fmt.Sprintf("%s/%s", config.UnixSocketLocation, config.UnixSocketName),
		mutex:      sync.Mutex{},
		metrics:    metrics,
		analyzer:   analyzer,
		logger:     logger,
	}
}

func (s *Socket) Start() error {
	os.Remove(s.socketPath)
	fmt.Println(s.socketPath)

	l, err := net.ListenUnix("unix", &net.UnixAddr{s.socketPath, "unix"})
	if err != nil {
		return fmt.Errorf("error listening on socket: %w", err)
		//s.logger.Error(fmt.Sprintf("Error listening on socket: %w", err))
	}

	s.listener = l

	socketListener = s.listener

	s.logger.Info(fmt.Sprintf("Listening on %s", s.socketPath))

	for isRunning {
		conn, err := l.Accept()
		if err != nil {
			if !isRunning {
				s.logger.Info("Socket server shutting down...")
				break
			}
			s.logger.Error(fmt.Sprintf("Error accepting connection: %s\n", err.Error()))
			continue
		}

		s.logger.Info(fmt.Sprintf("New client connected: %s", conn.RemoteAddr().String()))

		s.handleConnection(conn)
	}

	return nil
}

func (s *Socket) handleConnection(conn net.Conn) {
	defer conn.Close()

	sizeBuf := make([]byte, 4)
	_, err := io.ReadFull(conn, sizeBuf)
	if err != nil {
		fmt.Printf("Error reading message size :%s\n", err.Error())
		return
	}

	messageSize := binary.LittleEndian.Uint32(sizeBuf)

	messageBuf := make([]byte, messageSize)
	_, err = io.ReadFull(conn, messageBuf)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error reading message: %s\n", err.Error()))
		return
	}

	protoMessage := &frrProto.Message{}
	err = proto.Unmarshal(messageBuf, protoMessage)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error unmarshaling message: %s\n", err.Error()))
		return
	}

	//fmt.Printf("Received message: Command=%s, Package %s\n", protoMessage.Service, protoMessage.Command)

	execMutex.Lock()
	defer execMutex.Unlock()

	// TODO: Implement logging

	protoResponse := s.processCommand(protoMessage)

	responseData, err := proto.Marshal(protoResponse)
	if err != nil {
		fmt.Printf("Error marshaling response: %s\n", err.Error())
		s.logger.Error(fmt.Sprintf("Error marshaling response: %s\n", err.Error()))
		return
	}

	responseSizeBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(responseSizeBuf, uint32(len(responseData)))

	//fmt.Printf("Server marshaled %d bytes: %v\n", len(responseData), responseData)
	//fmt.Printf("Server buf size %d raw: %v\n", len(responseSizeBuf), responseSizeBuf)
	//fmt.Println(responseData)

	_, err = conn.Write(responseSizeBuf)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error sending response size: %s\n", err.Error()))
		return
	}

	_, err = conn.Write(responseData)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error sending response: %s\n", err.Error()))
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

func (s *Socket) exitSocketServer() {
	s.logger.Info("Shutting down socket server...")
	isRunning = false

	if socketListener != nil {
		socketListener.Close()
	}

	s.logger.Info("Socket server shut down completed")
	os.Exit(0)
}
