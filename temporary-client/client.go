package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	// Command line flags
	pkg := flag.String("package", "bgp", "Package to send command to (bgp, ospf, exit)")
	action := flag.String("action", "status", "Action to perform")
	socketPath := flag.String("socket", "/tmp/unixsock.sock", "Path to the Unix socket")
	flag.Parse()

	// Create the command
	cmd := map[string]interface{}{
		"package": *pkg,
		"action":  *action,
	}

	// Add parameters if any were provided as additional arguments
	params := make(map[string]interface{})
	args := flag.Args()
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			// Try to parse as number if possible
			var value interface{} = args[i+1]
			if num, err := json.Number(args[i+1]).Float64(); err == nil {
				value = num
			}
			params[args[i]] = value
		}
	}

	// Only add params if there are any
	if len(params) > 0 {
		cmd["params"] = params
	}

	// Convert to JSON
	jsonData, err := json.Marshal(cmd)
	if err != nil {
		fmt.Printf("Error creating JSON: %s\n", err.Error())
		os.Exit(1)
	}

	// Connect to the socket
	conn, err := net.Dial("unix", *socketPath)
	if err != nil {
		fmt.Printf("Error connecting to socket: %s\n", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	// Send the command
	_, err = conn.Write(jsonData)
	if err != nil {
		fmt.Printf("Error sending command: %s\n", err.Error())
		os.Exit(1)
	}

	// Read the response
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Error reading response: %s\n", err.Error())
		os.Exit(1)
	}

	// Pretty print the response
	var response map[string]interface{}
	if err := json.Unmarshal(buf[:n], &response); err != nil {
		fmt.Printf("Error parsing response: %s\n", err.Error())
		fmt.Printf("Raw response: %s\n", string(buf[:n]))
		os.Exit(1)
	}

	prettyJSON, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(prettyJSON))
}
