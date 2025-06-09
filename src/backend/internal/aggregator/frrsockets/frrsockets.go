package frrsockets

import (
	"bytes"
	"net"
	"path/filepath"
	"time"
)

// Package frrsockets provides FRR routing daemon communication via Unix sockets.
//
// This implementation is derived from the original work in tynany/frr_exporter:
// https://github.com/tynany/frr_exporter
//
// Original License:
// MIT License

type FRRCommandExecutor struct {
	DirPath string
	Timeout time.Duration
}

func NewConnection(dirPath string, timeout time.Duration) *FRRCommandExecutor {
	return &FRRCommandExecutor{DirPath: dirPath, Timeout: timeout}
}

func (c FRRCommandExecutor) ExecOSPFCmd(cmd string) ([]byte, error) {
	return executeCmd(filepath.Join(c.DirPath, "ospfd.vty"), cmd, c.Timeout)
}

func (c FRRCommandExecutor) ExecZebraCmd(cmd string) ([]byte, error) {
	return executeCmd(filepath.Join(c.DirPath, "zebra.vty"), cmd, c.Timeout)
}

func executeCmd(socketPath, cmd string, timeout time.Duration) ([]byte, error) {
	var response bytes.Buffer

	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{Net: "unix", Name: socketPath})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err = conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}

	buf := make([]byte, 4096)

	// Mimic vtysh by switching to 'enable' mode first. Note that commands need to be
	// null-terminated.
	if _, err = conn.Write([]byte("enable\x00")); err != nil {
		return nil, err
	}
	if _, err := conn.Read(buf); err != nil {
		return nil, err
	}

	// Send desired command.
	if _, err = conn.Write([]byte(cmd + "\x00")); err != nil {
		return nil, err
	}

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return response.Bytes(), err
		}

		response.Write(buf[:n])

		// frr signals the end of a response with a null character
		if n > 0 && buf[n-1] == 0 {
			return bytes.TrimRight(response.Bytes(), "\x00"), nil
		}
	}
}
