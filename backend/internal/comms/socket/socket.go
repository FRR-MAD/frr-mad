package socket

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	frrProto "github.com/ba2025-ysmprc/frr-tui/backend/proto"
	"google.golang.org/protobuf/proto"
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
		fmt.Printf("Error reading message: %s\n", err.Error())
		return
	}

	protoMessage := &frrProto.Message{}
	err = proto.Unmarshal(messageBuf, protoMessage)
	if err != nil {
		fmt.Printf("Error unmarshaling message: %s\n", err.Error())
		return
	}

	fmt.Printf("Received message: Command=%s, Package%s\n", protoMessage.Command, protoMessage.Package)

	execMutex.Lock()
	defer execMutex.Unlock()

	protoResponse := processCommand(protoMessage)

	responseData, err := proto.Marshal(protoResponse)
	if err != nil {
		fmt.Printf("Error marshaling response: %s\n", err.Error())
		return
	}

	responseSizeBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(responseSizeBuf, uint32(len(responseData)))

	fmt.Printf("Server marshaled %d bytes: %v\n", len(responseData), responseData)

	fmt.Printf("Server buf size %d raw: %v\n", len(responseSizeBuf), responseSizeBuf)

	fmt.Println(responseData)

	_, err = conn.Write(responseSizeBuf)
	if err != nil {
		fmt.Printf("Error sending response size: %s\n", err.Error())
		return
	}

	_, err = conn.Write(responseData)
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
