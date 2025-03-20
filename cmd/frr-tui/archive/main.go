package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	api "github.com/osrg/gobgp/v3/api"
	"github.com/osrg/gobgp/v3/pkg/server"
	"github.com/pterm/pterm"
	"google.golang.org/protobuf/types/known/anypb"
)

var testMode = flag.Bool("test", false, "Run in test mode to verify route anomaly detection")

// Enable verbose logging
var logger = log.New(os.Stdout, "[BGP-DEBUG] ", log.LstdFlags)

type RouteStatus struct {
	Status      string
	Description string
	Severity    string // "info", "warning", "error"
}

func main() {
	flag.Parse()
	logger.Println("Starting BGP server with debugging enabled")

	// Initialize GoBGP server
	bgpServer := server.NewBgpServer()
	go bgpServer.Serve()

	// Start the BGP server
	logger.Println("Configuring BGP global settings")
	if err := bgpServer.StartBgp(context.Background(), &api.StartBgpRequest{
		Global: &api.Global{
			Asn:        65001,
			RouterId:   "1.1.1.1",
			ListenPort: -1, // Let GoBGP choose the port
		},
	}); err != nil {
		logger.Fatalf("Failed to start BGP server: %v", err)
	}
	logger.Println("BGP server started successfully")

	expectedRoutes := map[string]string{
		"10.0.0.0/24":    "192.168.1.1",
		"172.16.0.0/16":  "192.168.1.2",
		"192.168.5.0/24": "192.168.1.3",
	}

	if *testMode {
		runTests(bgpServer, expectedRoutes)
	} else {
		addSimpleRoute(bgpServer)
		startTUI(bgpServer)
	}
}

func addSimpleRoute(bgpServer *server.BgpServer) {
	ctx := context.Background()
	logger.Println("Attempting to add a simple route")

	// Create a simple route: 10.0.0.0/24 via 192.168.1.1
	// Create NLRI
	nlri := &api.IPAddressPrefix{
		Prefix:    "10.0.0.0",
		PrefixLen: 24,
	}
	nlriAny, err := anypb.New(nlri)
	if err != nil {
		logger.Printf("Failed to marshal NLRI: %v", err)
		return
	}
	logger.Printf("Created NLRI: %s/%d", nlri.Prefix, nlri.PrefixLen)

	// Create path attributes
	origin := &api.OriginAttribute{
		Origin: 0, // IGP
	}
	originAny, err := anypb.New(origin)
	if err != nil {
		logger.Printf("Failed to marshal origin: %v", err)
		return
	}
	logger.Printf("Created origin attribute: %d", origin.Origin)

	nextHop := &api.NextHopAttribute{
		NextHop: "192.168.1.1",
	}
	nextHopAny, err := anypb.New(nextHop)
	if err != nil {
		logger.Printf("Failed to marshal next-hop: %v", err)
		return
	}
	logger.Printf("Created next-hop attribute: %s", nextHop.NextHop)

	// Create path
	path := &api.Path{
		Family: &api.Family{
			Afi:  api.Family_AFI_IP,
			Safi: api.Family_SAFI_UNICAST,
		},
		Nlri:   nlriAny,
		Pattrs: []*anypb.Any{originAny, nextHopAny},
	}

	// Add path to BGP server
	logger.Printf("Adding path to BGP server, TypeUrl: %s", nlriAny.TypeUrl)
	for i, attr := range path.Pattrs {
		logger.Printf("Attribute[%d] TypeUrl: %s", i, attr.TypeUrl)
	}

	_, err = bgpServer.AddPath(ctx, &api.AddPathRequest{
		TableType: api.TableType_GLOBAL,
		Path:      path,
	})
	if err != nil {
		logger.Printf("Failed to add path: %v", err)

		// Check if the error is related to the NLRI format
		if strings.Contains(err.Error(), "nlri") {
			logger.Println("NLRI format issue detected. Trying alternative formatting...")

			// Try alternative formatting
			tryAlternativeFormats(bgpServer)
		}
	} else {
		logger.Println("Path added successfully!")
	}
}

func tryAlternativeFormats(bgpServer *server.BgpServer) {
	ctx := context.Background()
	logger.Println("Trying alternative NLRI formats")

	// Try with a different prefix format
	prefixes := []string{"10.0.0.0/24", "10.0.0.0", "10.0.0"}
	for i, prefix := range prefixes {
		logger.Printf("Attempt %d: Trying with prefix format: %s", i+1, prefix)

		// Parse prefix
		_, ipNet, err := net.ParseCIDR(prefix)
		if err != nil {
			// If not CIDR, try as IP address
			ip := net.ParseIP(prefix)
			if ip == nil {
				logger.Printf("Invalid prefix format: %s", prefix)
				continue
			}
			// Create a /24 subnet
			ipNet = &net.IPNet{
				IP:   ip,
				Mask: net.CIDRMask(24, 32),
			}
		}

		// Create NLRI
		nlri := &api.IPAddressPrefix{
			Prefix:    ipNet.IP.String(),
			PrefixLen: 24,
		}
		nlriAny, err := anypb.New(nlri)
		if err != nil {
			logger.Printf("Failed to marshal NLRI: %v", err)
			continue
		}

		// Create attributes
		origin := &api.OriginAttribute{Origin: 0}
		originAny, err := anypb.New(origin)
		if err != nil {
			logger.Printf("Failed to marshal origin: %v", err)
			continue
		}

		nextHop := &api.NextHopAttribute{NextHop: "192.168.1.1"}
		nextHopAny, err := anypb.New(nextHop)
		if err != nil {
			logger.Printf("Failed to marshal next-hop: %v", err)
			continue
		}

		// Add path
		_, err = bgpServer.AddPath(ctx, &api.AddPathRequest{
			TableType: api.TableType_GLOBAL,
			Path: &api.Path{
				Family: &api.Family{
					Afi:  api.Family_AFI_IP,
					Safi: api.Family_SAFI_UNICAST,
				},
				Nlri:   nlriAny,
				Pattrs: []*anypb.Any{originAny, nextHopAny},
			},
		})

		if err != nil {
			logger.Printf("Attempt %d failed: %v", i+1, err)
		} else {
			logger.Printf("Attempt %d succeeded!", i+1)
			break
		}
	}
}

func startTUI(bgpServer *server.BgpServer) {
	logger.Println("Starting BGP monitoring TUI")
	pterm.Info.Println("Starting BGP Monitoring TUI with debugging...")

	// Define expected routes (you can make this configurable)
	expectedRoutes := map[string]string{
		"10.0.0.0/24": "192.168.1.1",
		// Add more expected routes here
	}

	// Channel to signal when to refresh the TUI
	refreshTUI := make(chan bool, 1)

	go func() {
		for {
			select {
			case <-refreshTUI:
				return
			default:
				pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)).WithMargin(10).Println("BGP Route Monitor")

				// Fetch BGP routes with error handling
				var routes []*api.Path
				err := bgpServer.ListPath(context.Background(), &api.ListPathRequest{
					TableType: api.TableType_GLOBAL,
					Family:    &api.Family{Afi: api.Family_AFI_IP, Safi: api.Family_SAFI_UNICAST},
				}, func(destination *api.Destination) {
					routes = append(routes, destination.Paths...)
					logger.Printf("Found destination: %s with %d paths", destination.Prefix, len(destination.Paths))
				})

				if err != nil {
					pterm.Error.Printf("Failed to fetch routes: %v\n", err)
					logger.Printf("ListPath error: %v", err)
				} else {
					logger.Printf("Found %d routes", len(routes))
				}

				// Validate routes
				routeStatus := validateRoutes(routes, expectedRoutes)

				// Display routes in TUI
				tableData := pterm.TableData{
					{"Prefix", "Next Hop", "Origin", "Status", "Alert"},
				}

				// Show warning if no routes found
				if len(routes) == 0 {
					pterm.Warning.Println("No routes found in BGP table!")
					tableData = append(tableData, []string{
						"N/A", "N/A", "N/A", "No Routes", "Check logs for errors",
					})
				}

				// Display actual routes
				for _, route := range routes {
					prefix := "Unknown"
					nextHop := "Unknown"
					origin := "Unknown"
					status := "ACTIVE"
					alert := ""

					// Get the prefix
					if route.GetNlri() != nil {
						if ipPrefix, err := getIPPrefixFromNLRI(route.GetNlri()); err == nil {
							prefix = ipPrefix
							// Get status from validation
							if routeStat, ok := routeStatus[prefix]; ok {
								status = routeStat.Status
								alert = routeStat.Description
							}
						} else {
							prefix = "Error: " + err.Error()
						}
					}

					// Extract attributes
					for _, attr := range route.GetPattrs() {
						// Extract next hop
						if strings.Contains(attr.TypeUrl, "NextHop") {
							nextHopAttr := &api.NextHopAttribute{}
							if err := attr.UnmarshalTo(nextHopAttr); err == nil {
								nextHop = nextHopAttr.GetNextHop()
							}
						}

						// Extract origin
						if strings.Contains(attr.TypeUrl, "Origin") {
							originAttr := &api.OriginAttribute{}
							if err := attr.UnmarshalTo(originAttr); err == nil {
								switch originAttr.GetOrigin() {
								case 0:
									origin = "IGP"
								case 1:
									origin = "EGP"
								case 2:
									origin = "INCOMPLETE"
								}
							}
						}
					}

					tableData = append(tableData, []string{
						prefix,
						nextHop,
						origin,
						status,
						alert,
					})
				}

				// Display missing routes
				for prefix, nextHop := range expectedRoutes {
					if status, ok := routeStatus[prefix]; ok && status.Status == "MISSING" {
						tableData = append(tableData, []string{
							prefix,
							nextHop,
							"N/A",
							"MISSING",
							status.Description,
						})
					}
				}

				// Render the table with colored status
				table := pterm.DefaultTable.WithHasHeader().WithData(tableData)
				table.Render()

				// Display alert summary
				pterm.Println("\nAlert Summary:")
				hasAlerts := false
				for prefix, status := range routeStatus {
					if status.Status != "ACTIVE" {
						hasAlerts = true
						switch status.Severity {
						case "error":
							pterm.Error.Printf("Route %s: %s - %s\n", prefix, status.Status, status.Description)
						case "warning":
							pterm.Warning.Printf("Route %s: %s - %s\n", prefix, status.Status, status.Description)
						default:
							pterm.Info.Printf("Route %s: %s - %s\n", prefix, status.Status, status.Description)
						}
					}
				}

				if !hasAlerts {
					pterm.Success.Println("No issues detected with current routes.")
				}

				pterm.Println("\nBGP Server Information:")
				pterm.Printf("- Routes found: %d\n", len(routes))
				pterm.Printf("- Expected routes: %d\n", len(expectedRoutes))
				pterm.Printf("- ASN: 65001\n")
				pterm.Printf("- Router ID: 1.1.1.1\n")
				pterm.Println("\nPress Ctrl+C to exit. Refreshing data every 5 seconds...")

				time.Sleep(5 * time.Second)
			}
		}
	}()

	// Goroutine to handle manual GoBGP CLI commands
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			pterm.Info.Print("Enter GoBGP CLI command: ")
			command, _ := reader.ReadString('\n')
			command = strings.TrimSpace(command)

			if command == "" {
				continue
			}

			// Execute the GoBGP CLI command
			cmd := exec.Command("sh", "-c", command)
			output, err := cmd.CombinedOutput()
			if err != nil {
				pterm.Error.Printf("Command failed: %v\n", err)
			}

			// Display the command output
			pterm.Println("\nCommand Output:")
			pterm.Println(string(output))
		}
	}()

	// Wait for Ctrl+C to exit
	<-make(chan struct{})
}

func runTests(bgpServer *server.BgpServer, expectedRoutes map[string]string) {
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)).WithMargin(10).Println("BGP Route Anomaly Tests")

	// Define test cases
	testCases := []struct {
		name        string
		description string
		setup       func(bgpServer *server.BgpServer)
		expectError string
	}{
		{
			name:        "Normal Routes",
			description: "All routes are correctly advertised",
			setup: func(bgpServer *server.BgpServer) {
				// Clear all routes
				clearAllRoutes(bgpServer)

				// Add all expected routes correctly
				for prefix, nextHop := range expectedRoutes {
					parts := strings.Split(prefix, "/")
					prefixIP := parts[0]
					prefixLen := 24 // Default to /24
					if len(parts) > 1 {
						fmt.Sscanf(parts[1], "%d", &prefixLen)
					}
					addRoute(bgpServer, prefixIP, uint32(prefixLen), nextHop)
				}
			},
			expectError: "",
		},
		{
			name:        "Missing Route",
			description: "One expected route is missing",
			setup: func(bgpServer *server.BgpServer) {
				// Clear all routes
				clearAllRoutes(bgpServer)

				// Add only some of the expected routes
				first := true
				for prefix, nextHop := range expectedRoutes {
					if first {
						// Skip the first route to simulate missing route
						first = false
						continue
					}

					parts := strings.Split(prefix, "/")
					prefixIP := parts[0]
					prefixLen := 24 // Default to /24
					if len(parts) > 1 {
						fmt.Sscanf(parts[1], "%d", &prefixLen)
					}
					addRoute(bgpServer, prefixIP, uint32(prefixLen), nextHop)
				}
			},
			expectError: "MISSING",
		},
		{
			name:        "Duplicate Route",
			description: "A route is advertised multiple times",
			setup: func(bgpServer *server.BgpServer) {
				// Clear all routes
				clearAllRoutes(bgpServer)

				// Add all expected routes
				for prefix, nextHop := range expectedRoutes {
					parts := strings.Split(prefix, "/")
					prefixIP := parts[0]
					prefixLen := 24 // Default to /24
					if len(parts) > 1 {
						fmt.Sscanf(parts[1], "%d", &prefixLen)
					}
					addRoute(bgpServer, prefixIP, uint32(prefixLen), nextHop)

					// Add the first route again to create a duplicate
					if prefix == "10.0.0.0/24" {
						addRoute(bgpServer, prefixIP, uint32(prefixLen), nextHop)
					}
				}
			},
			expectError: "DUPLICATE",
		},
		{
			name:        "Wrong Next Hop",
			description: "A route has incorrect next hop",
			setup: func(bgpServer *server.BgpServer) {
				// Clear all routes
				clearAllRoutes(bgpServer)

				// Add routes, but with one having wrong next hop
				for prefix, nextHop := range expectedRoutes {
					parts := strings.Split(prefix, "/")
					prefixIP := parts[0]
					prefixLen := 24 // Default to /24
					if len(parts) > 1 {
						fmt.Sscanf(parts[1], "%d", &prefixLen)
					}

					if prefix == "172.16.0.0/16" {
						// Use wrong next hop for this route
						addRoute(bgpServer, prefixIP, uint32(prefixLen), "192.168.2.1")
					} else {
						addRoute(bgpServer, prefixIP, uint32(prefixLen), nextHop)
					}
				}
			},
			expectError: "WRONG_NEXTHOP",
		},
		{
			name:        "Unexpected Route",
			description: "A route not in expected list is advertised",
			setup: func(bgpServer *server.BgpServer) {
				// Clear all routes
				clearAllRoutes(bgpServer)

				// Add all expected routes
				for prefix, nextHop := range expectedRoutes {
					parts := strings.Split(prefix, "/")
					prefixIP := parts[0]
					prefixLen := 24 // Default to /24
					if len(parts) > 1 {
						fmt.Sscanf(parts[1], "%d", &prefixLen)
					}
					addRoute(bgpServer, prefixIP, uint32(prefixLen), nextHop)
				}

				// Add an unexpected route
				addRoute(bgpServer, "8.8.8.0", 24, "192.168.1.10")
			},
			expectError: "UNEXPECTED",
		},
	}

	// Run each test
	for _, tc := range testCases {
		pterm.DefaultSection.Println(tc.name)
		pterm.Println(tc.description)

		// Set up the routes for this test
		tc.setup(bgpServer)

		// Wait for routes to propagate
		time.Sleep(500 * time.Millisecond)

		// Fetch and validate routes
		routes := fetchRoutes(bgpServer)
		routeStatus := validateRoutes(routes, expectedRoutes)

		// Check if the expected error is found
		testPassed := false
		if tc.expectError == "" {
			// Expect no errors
			testPassed = true
			for _, status := range routeStatus {
				if status.Status != "ACTIVE" {
					testPassed = false
					break
				}
			}
		} else {
			// Expect specific error
			for _, status := range routeStatus {
				if status.Status == tc.expectError {
					testPassed = true
					break
				}
			}
		}

		// Display routes and alerts
		displayRoutes(routes, routeStatus)
		displayAlerts(routeStatus)

		// Display test result
		if testPassed {
			pterm.Success.Println("Test PASSED: Expected condition detected")
		} else {
			pterm.Error.Println("Test FAILED: Expected condition not detected")
		}

		pterm.Println("\nPress Enter to continue to next test...")
		fmt.Scanln()
	}

	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).WithMargin(10).Println("All Tests Completed")
}

// Helper functions

func addRoute(bgpServer *server.BgpServer, prefix string, prefixLen uint32, nextHop string) {
	ctx := context.Background()
	logger.Printf("Adding route %s/%d via %s", prefix, prefixLen, nextHop)

	// Create NLRI
	nlri := &api.IPAddressPrefix{
		Prefix:    prefix,
		PrefixLen: prefixLen,
	}
	nlriAny, err := anypb.New(nlri)
	if err != nil {
		logger.Printf("Failed to marshal NLRI: %v", err)
		return
	}

	// Create attributes
	origin := &api.OriginAttribute{Origin: 0} // IGP
	originAny, err := anypb.New(origin)
	if err != nil {
		logger.Printf("Failed to marshal origin: %v", err)
		return
	}

	nextHopAttr := &api.NextHopAttribute{NextHop: nextHop}
	nextHopAny, err := anypb.New(nextHopAttr)
	if err != nil {
		logger.Printf("Failed to marshal next-hop: %v", err)
		return
	}

	// Add path
	_, err = bgpServer.AddPath(ctx, &api.AddPathRequest{
		TableType: api.TableType_GLOBAL,
		Path: &api.Path{
			Family: &api.Family{
				Afi:  api.Family_AFI_IP,
				Safi: api.Family_SAFI_UNICAST,
			},
			Nlri:   nlriAny,
			Pattrs: []*anypb.Any{originAny, nextHopAny},
		},
	})

	if err != nil {
		logger.Printf("Failed to add path: %v", err)
	} else {
		logger.Printf("Path added successfully!")
	}
}

func clearAllRoutes(bgpServer *server.BgpServer) {
	ctx := context.Background()
	logger.Println("Clearing all routes")

	var routes []*api.Path
	err := bgpServer.ListPath(ctx, &api.ListPathRequest{
		TableType: api.TableType_GLOBAL,
		Family:    &api.Family{Afi: api.Family_AFI_IP, Safi: api.Family_SAFI_UNICAST},
	}, func(destination *api.Destination) {
		routes = append(routes, destination.Paths...)
	})

	if err != nil {
		logger.Printf("Failed to list paths: %v", err)
		return
	}

	for _, route := range routes {
		result := bgpServer.DeletePath(ctx, &api.DeletePathRequest{
			TableType: api.TableType_GLOBAL,
			Path:      route,
		})
		if result != nil {
			logger.Printf("Failed to delete path: %v", err)
		}
	}

	logger.Printf("Cleared %d routes", len(routes))
}

func fetchRoutes(bgpServer *server.BgpServer) []*api.Path {
	var routes []*api.Path
	err := bgpServer.ListPath(context.Background(), &api.ListPathRequest{
		TableType: api.TableType_GLOBAL,
		Family:    &api.Family{Afi: api.Family_AFI_IP, Safi: api.Family_SAFI_UNICAST},
	}, func(destination *api.Destination) {
		routes = append(routes, destination.Paths...)
		logger.Printf("Found destination: %s with %d paths", destination.Prefix, len(destination.Paths))
	})

	if err != nil {
		logger.Printf("ListPath error: %v", err)
	} else {
		logger.Printf("Found %d routes", len(routes))
	}

	return routes
}

// ValidateRoutes checks for route issues
func validateRoutes(routes []*api.Path, expectedRoutes map[string]string) map[string]RouteStatus {
	routeStatus := make(map[string]RouteStatus)
	foundPrefixes := make(map[string]int)

	// Count occurrences of each prefix
	for _, route := range routes {
		prefix, err := getIPPrefixFromNLRI(route.GetNlri())
		if err != nil {
			continue
		}
		foundPrefixes[prefix]++
	}

	// Check for duplicates
	for prefix, count := range foundPrefixes {
		if count > 1 {
			routeStatus[prefix] = RouteStatus{
				Status:      "DUPLICATE",
				Description: fmt.Sprintf("Route advertised %d times", count),
				Severity:    "error",
			}
		} else {
			routeStatus[prefix] = RouteStatus{
				Status:      "ACTIVE",
				Description: "Route is correctly advertised",
				Severity:    "info",
			}
		}
	}

	// Check for wrongly advertised routes (unexpected routes)
	for prefix := range foundPrefixes {
		if _, ok := expectedRoutes[prefix]; !ok {
			routeStatus[prefix] = RouteStatus{
				Status:      "UNEXPECTED",
				Description: "Route is not in the expected routes list",
				Severity:    "warning",
			}
		}
	}

	// Check for unadvertised routes (missing routes)
	for prefix, nextHop := range expectedRoutes {
		if _, ok := foundPrefixes[prefix]; !ok {
			routeStatus[prefix] = RouteStatus{
				Status:      "MISSING",
				Description: fmt.Sprintf("Expected route %s via %s is not advertised", prefix, nextHop),
				Severity:    "error",
			}
		} else {
			// Verify next hop if the route exists
			for _, route := range routes {
				prefixStr, err := getIPPrefixFromNLRI(route.GetNlri())
				if err != nil || prefixStr != prefix {
					continue
				}

				routeNextHop := getNextHopFromPath(route)
				if routeNextHop != nextHop {
					routeStatus[prefix] = RouteStatus{
						Status:      "WRONG_NEXTHOP",
						Description: fmt.Sprintf("Expected next hop %s, got %s", nextHop, routeNextHop),
						Severity:    "warning",
					}
				}
			}
		}
	}

	return routeStatus
}

func displayRoutes(routes []*api.Path, routeStatus map[string]RouteStatus) {
	// Display routes in TUI
	tableData := pterm.TableData{
		{"Prefix", "Next Hop", "Origin", "Status", "Alert"},
	}

	// Show warning if no routes found
	if len(routes) == 0 {
		pterm.Warning.Println("No routes found in BGP table!")
		tableData = append(tableData, []string{
			"N/A", "N/A", "N/A", "No Routes", "Check logs for errors",
		})
	}

	// Display actual routes
	for _, route := range routes {
		prefix := "Unknown"
		nextHop := "Unknown"
		origin := "Unknown"
		status := "ACTIVE"
		alert := ""

		// Get the prefix
		if route.GetNlri() != nil {
			if ipPrefix, err := getIPPrefixFromNLRI(route.GetNlri()); err == nil {
				prefix = ipPrefix
				// Get status from validation
				if routeStat, ok := routeStatus[prefix]; ok {
					status = routeStat.Status
					alert = routeStat.Description
				}
			} else {
				prefix = "Error: " + err.Error()
			}
		}

		// Extract attributes
		for _, attr := range route.GetPattrs() {
			// Extract next hop
			if strings.Contains(attr.TypeUrl, "NextHop") {
				nextHopAttr := &api.NextHopAttribute{}
				if err := attr.UnmarshalTo(nextHopAttr); err == nil {
					nextHop = nextHopAttr.GetNextHop()
				}
			}

			// Extract origin
			if strings.Contains(attr.TypeUrl, "Origin") {
				originAttr := &api.OriginAttribute{}
				if err := attr.UnmarshalTo(originAttr); err == nil {
					switch originAttr.GetOrigin() {
					case 0:
						origin = "IGP"
					case 1:
						origin = "EGP"
					case 2:
						origin = "INCOMPLETE"
					}
				}
			}
		}

		tableData = append(tableData, []string{
			prefix,
			nextHop,
			origin,
			status,
			alert,
		})
	}

	// Render the table
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func displayAlerts(routeStatus map[string]RouteStatus) {
	pterm.Println("\nAlert Summary:")
	hasAlerts := false
	for prefix, status := range routeStatus {
		if status.Status != "ACTIVE" {
			hasAlerts = true
			switch status.Severity {
			case "error":
				pterm.Error.Printf("Route %s: %s - %s\n", prefix, status.Status, status.Description)
			case "warning":
				pterm.Warning.Printf("Route %s: %s - %s\n", prefix, status.Status, status.Description)
			default:
				pterm.Info.Printf("Route %s: %s - %s\n", prefix, status.Status, status.Description)
			}
		}
	}

	if !hasAlerts {
		pterm.Success.Println("No issues detected with current routes.")
	}
}

// Helper function to get next hop from path
func getNextHopFromPath(path *api.Path) string {
	for _, attr := range path.GetPattrs() {
		if strings.Contains(attr.TypeUrl, "NextHop") {
			nextHopAttr := &api.NextHopAttribute{}
			if err := attr.UnmarshalTo(nextHopAttr); err == nil {
				return nextHopAttr.GetNextHop()
			}
		}
	}
	return "Unknown"
}

func getIPPrefixFromNLRI(nlri *anypb.Any) (string, error) {
	if strings.Contains(nlri.TypeUrl, "IPAddressPrefix") {
		ipPrefix := &api.IPAddressPrefix{}
		if err := nlri.UnmarshalTo(ipPrefix); err != nil {
			return "", err
		}
		return fmt.Sprintf("%s/%d", ipPrefix.GetPrefix(), ipPrefix.GetPrefixLen()), nil
	}
	return "", fmt.Errorf("unsupported NLRI type: %s", nlri.TypeUrl)
}

// not used yet
func fetchRoutesWithCLI() []string {
	logger.Println("Fetching routes using GoBGP CLI...")
	cmd := exec.Command("gobgp", "global", "rib")
	output, err := cmd.Output()
	if err != nil {
		logger.Printf("Failed to fetch routes using GoBGP CLI: %v", err)
		return nil
	}
	// Parse CLI output
	routes := strings.Split(string(output), "\n")
	var filteredRoutes []string
	for _, route := range routes {
		if strings.Contains(route, "/") {
			filteredRoutes = append(filteredRoutes, route)
		}
	}
	return filteredRoutes
}
