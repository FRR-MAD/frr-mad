package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	api "github.com/osrg/gobgp/v3/api"
	"github.com/osrg/gobgp/v3/pkg/server"
	"github.com/pterm/pterm"
	"google.golang.org/protobuf/types/known/anypb"
)

// Enable verbose logging
var logger = log.New(os.Stdout, "[BGP-DEBUG] ", log.LstdFlags)

func main() {
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

	// Try adding a route with the direct method
	addSimpleRoute(bgpServer)

	// Start TUI
	startTUI(bgpServer)
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

	for {
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

		// Display routes in TUI
		tableData := pterm.TableData{
			{"Prefix", "Next Hop", "Origin", "Status", "Debug Info"},
		}

		// Show warning if no routes found
		if len(routes) == 0 {
			pterm.Warning.Println("No routes found in BGP table!")
			tableData = append(tableData, []string{
				"N/A", "N/A", "N/A", "No Routes", "Check logs for errors",
			})
		}

		for i, route := range routes {
			prefix := "Unknown"
			nextHop := "Unknown"
			origin := "Unknown"
			status := "Active"

			// Get the prefix from the destination
			if route.GetNlri() != nil {
				if ipPrefix, err := getIPPrefixFromNLRI(route.GetNlri()); err == nil {
					prefix = ipPrefix
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
				fmt.Sprintf("Path[%d]: Age: %ds", i, route.GetAge().GetSeconds()),
			})
		}

		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

		// Display debug info
		pterm.Println("\nBGP Server Information:")
		pterm.Printf("- Routes found: %d\n", len(routes))
		pterm.Printf("- ASN: 65001\n")
		pterm.Printf("- Router ID: 1.1.1.1\n")
		pterm.Println("\nPress Ctrl+C to exit. Refreshing data every 5 seconds...")

		time.Sleep(5 * time.Second)
	}
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
