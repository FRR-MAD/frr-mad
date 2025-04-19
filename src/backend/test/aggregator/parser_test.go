package aggregator_test

import (
	"fmt"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"os"
	"path/filepath"
	"testing"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/aggregator"
)

func TestParseStaticFRRConfig(t *testing.T) {
	// Test valid config
	configPath := "./mock-files/r101.conf"
	config, err := aggregator.ParseStaticFRRConfig(configPath)
	if err != nil {
		t.Fatalf("ParseStaticFRRConfig failed: %v", err)
	}

	fmt.Printf("ParseStaticFRRConfig: %v\n", config)

	// v=========== METADATA TEST ==========v
	if config.Hostname != "r101" {
		t.Errorf("Expected Hostname to be 'r101', got '%s'", config.Hostname)
	}
	if config.FrrVersion != "8.5.4_git" {
		t.Errorf("Expected FrrVersion to be '8.5.4_git', got '%s'", config.FrrVersion)
	}
	if config.ServiceAdvancedVty != true {
		t.Error("Expected ServiceAdvancedVty to be true")
	}

	if len(config.Interfaces) != 12 {
		t.Errorf("Expected 13 interfaces, got %d", len(config.Interfaces))
	}
	// ^=========== METADATA TEST ==========^

	// v=========== INTERFACE TEST ==========v
	eth1 := config.Interfaces[0]
	if eth1.Name != "eth1" {
		t.Errorf("Expected first interface name to be 'eth1', got '%s'", eth1.Name)
	}
	if len(eth1.IpAddress) == 0 {
		t.Error("Expected eth1 to have at least one IP address")
	} else {
		if eth1.IpAddress[0].IpAddress != "172.22.1.1" {
			t.Errorf("Expected eth1 IP to be 172.22.1.1, got '%s'", eth1.IpAddress[0].IpAddress)
		}
		if eth1.IpAddress[0].PrefixLength != 24 {
			t.Errorf("Expected eth1 prefix length to be 24, got %d", eth1.IpAddress[0].PrefixLength)
		}
	}

	// eth2
	eth2 := config.Interfaces[1]
	if eth2.Name != "eth2" {
		t.Errorf("Expected second interface name to be 'eth2', got '%s'", eth2.Name)
	}
	if len(eth2.IpAddress) == 0 {
		t.Error("Expected eth2 to have at least one IP address")
	} else {
		if eth2.IpAddress[0].IpAddress != "10.0.12.1" {
			t.Errorf("Expected eth2 IP to be 172.22.1.1, got '%s'", eth2.IpAddress[0].IpAddress)
		}
		if eth2.IpAddress[0].PrefixLength != 24 {
			t.Errorf("Expected eth2 prefix length to be 24, got %d", eth2.IpAddress[0].PrefixLength)
		}
	}
	if eth2.Area != "0.0.0.0" {
		t.Errorf("Expected eth2 to be in Area '0.0.0.0', got '%s'", eth2.Area)
	}
	if eth2.Passive {
		t.Errorf("Expected eth2 to be in Passive 'false', got '%v'", eth2.Passive)
	}

	// eth3
	eth3 := config.Interfaces[2]
	if eth3.Name != "eth3" {
		t.Errorf("Expected eth3 name to be 'eth3', got '%s'", eth3.Name)
	}
	if len(eth3.IpAddress) == 0 {
		t.Error("Expected eth3 to have at least one IP address")
	} else {
		if eth3.IpAddress[0].IpAddress != "10.0.13.1" {
			t.Errorf("Expected eth3 IP to be 10.0.13.1, got '%s'", eth3.IpAddress[0].IpAddress)
		}
		if eth3.IpAddress[0].PrefixLength != 24 {
			t.Errorf("Expected eth3 prefix length to be 24, got %d", eth3.IpAddress[0].PrefixLength)
		}
	}
	if eth3.Area != "0.0.0.0" {
		t.Errorf("Expected eth3 to be in Area '0.0.0.0', got '%s'", eth3.Area)
	}
	if eth3.Passive {
		t.Error("Expected eth3 to be non-passive")
	}

	// eth4
	eth4 := config.Interfaces[3]
	if eth4.Name != "eth4" {
		t.Errorf("Expected eth4 name to be 'eth4', got '%s'", eth4.Name)
	}
	if len(eth4.IpAddress) == 0 {
		t.Error("Expected eth4 to have at least one IP address")
	} else {
		if eth4.IpAddress[0].IpAddress != "10.0.0.1" {
			t.Errorf("Expected eth4 IP to be 10.0.0.1, got '%s'", eth4.IpAddress[0].IpAddress)
		}
		if eth4.IpAddress[0].PrefixLength != 23 {
			t.Errorf("Expected eth4 prefix length to be 23, got %d", eth4.IpAddress[0].PrefixLength)
		}
	}
	if eth4.Area != "0.0.0.0" {
		t.Errorf("Expected eth4 to be in Area '0.0.0.0', got '%s'", eth4.Area)
	}
	if !eth4.Passive {
		t.Error("Expected eth4 to be passive")
	}

	// eth11
	eth11 := config.Interfaces[10]
	if eth11.Name != "eth11" {
		t.Errorf("Expected eth11 name to be 'eth11', got '%s'", eth11.Name)
	}
	if len(eth11.IpAddress) == 0 {
		t.Error("Expected eth11 to have at least one IP address")
	} else {
		if eth11.IpAddress[0].IpAddress != "10.0.19.1" {
			t.Errorf("Expected eth11 IP to be 10.0.19.1, got '%s'", eth11.IpAddress[0].IpAddress)
		}
		if eth11.IpAddress[0].PrefixLength != 24 {
			t.Errorf("Expected eth11 prefix length to be 24, got %d", eth11.IpAddress[0].PrefixLength)
		}
	}
	if eth11.Area != "0.0.0.0" {
		t.Errorf("Expected eth11 to be in Area '0.0.0.0', got '%s'", eth11.Area)
	}
	if eth11.Passive {
		t.Error("Expected eth11 to be non-passive")
	}

	// lo
	lo := config.Interfaces[11]
	if lo.Name != "lo" {
		t.Errorf("Expected lo interface name to be 'lo', got '%s'", lo.Name)
	}
	if len(lo.IpAddress) == 0 {
		t.Error("Expected lo to have at least one IP address")
	} else {
		if lo.IpAddress[0].IpAddress != "65.0.1.1" {
			t.Errorf("Expected lo IP to be 65.0.1.1, got '%s'", lo.IpAddress[0].IpAddress)
		}
		if lo.IpAddress[0].PrefixLength != 32 {
			t.Errorf("Expected lo prefix length to be 32, got %d", lo.IpAddress[0].PrefixLength)
		}
	}
	if !lo.Passive {
		t.Error("Expected lo to be passive")
	}
	// ^=========== INTERFACE TEST ==========^

	// v=========== STATIC ROUTE TEST ==========v
	staticRoutes := config.StaticRoutes
	if len(staticRoutes) != 1 {
		t.Errorf("Expected 1 static route, got %d", len(staticRoutes))
	}
	if staticRoutes[0].IpPrefix.IpAddress != "192.168.1.0" {
		t.Errorf("Expected static route one to have 192.168.1.0, got '%s'", staticRoutes[0].IpPrefix)
	}
	if staticRoutes[0].IpPrefix.PrefixLength != 24 {
		t.Errorf("Expected static route Prefix one to be 24, got %d", staticRoutes[0].IpPrefix.PrefixLength)
	}
	if staticRoutes[0].NextHop != "192.168.100.91" {
		t.Errorf("Expected Next Hop to be 192.168.100.91, got '%s'", staticRoutes[0].NextHop)
	}
	// ^=========== STATIC ROUTE TEST ==========^

	// v=========== OSPF CONFIG TEST ==========v
	ospfConfig := config.OspfConfig
	if ospfConfig.RouterId != "65.0.1.1" {
		t.Errorf("Expected ospf router id 65.0.1.1, got '%s'", ospfConfig.RouterId)
	}
	if ospfConfig.Redistribution[0].Type != "static" {
		t.Errorf("Expected ospf redistribution type static, got '%s'", ospfConfig.Redistribution[0].Type)
	}
	if ospfConfig.Redistribution[0].Metric != "1" {
		t.Errorf("Expected Metric to be '1', got '%s'", ospfConfig.Redistribution[0].Metric)
	}
	if ospfConfig.Redistribution[0].RouteMap != "lanroutes" {
		t.Errorf("Expected RouteMap to be 'lanroutes', got '%s'", ospfConfig.Redistribution[0].RouteMap)
	}
	if ospfConfig.Redistribution[1].Type != "bgp" {
		t.Errorf("Expected Redistribution type bgp, got '%s'", ospfConfig.Redistribution[1].Type)
	}
	if ospfConfig.Redistribution[1].Metric != "1" {
		t.Errorf("Expected Metric to be '1', got '%s'", ospfConfig.Redistribution[1].Metric)
	}
	if ospfConfig.Redistribution[1].RouteMap != "" {
		t.Errorf("Expected RouteMap to be empty, got '%s'", ospfConfig.Redistribution[1].RouteMap)
	}
	// ^=========== OSPF CONFIG TEST ==========^

	// v=========== ACCESS LIST TEST ==========v
	accessList := config.AccessList
	if len(accessList) != 2 {
		t.Errorf("Expected 2 access list, got %d", len(accessList))
	}
	termList, ok := config.AccessList["term"]
	if !ok {
		t.Fatalf("Access list 'term' not found")
	}
	if len(termList.AccessListItems) != 2 {
		t.Errorf("Expected 2 items in access-list 'term', got %d", len(termList.AccessListItems))
	}

	termItem1 := termList.AccessListItems[0]
	if termItem1.Sequence != 5 || termItem1.AccessControl != "permit" {
		t.Errorf("Unexpected access-list term item 1: %+v", termItem1)
	}
	if prefix, ok := termItem1.Destination.(*frrProto.AccessListItem_IpPrefix); !ok {
		t.Errorf("Expected term item 1 to be IP prefix, got: %+v", termItem1.Destination)
	} else {
		if prefix.IpPrefix.IpAddress != "127.0.0.1" || prefix.IpPrefix.PrefixLength != 32 {
			t.Errorf("Expected 127.0.0.1/32, got %s/%d", prefix.IpPrefix.IpAddress, prefix.IpPrefix.PrefixLength)
		}
	}

	termItem2 := termList.AccessListItems[1]
	if termItem2.Sequence != 10 || termItem2.AccessControl != "deny" {
		t.Errorf("Unexpected access-list term item 2: %+v", termItem2)
	}
	if _, ok := termItem2.Destination.(*frrProto.AccessListItem_Any); !ok {
		t.Errorf("Expected term item 2 to be 'any', got: %+v", termItem2.Destination)
	}

	localsiteList, ok := config.AccessList["localsite"]
	if !ok {
		t.Fatalf("Access list 'localsite' not found")
	}
	if len(localsiteList.AccessListItems) != 1 {
		t.Errorf("Expected 1 item in access-list 'localsite', got %d", len(localsiteList.AccessListItems))
	}

	localsiteItem := localsiteList.AccessListItems[0]
	if localsiteItem.Sequence != 15 || localsiteItem.AccessControl != "permit" {
		t.Errorf("Unexpected localsite item: %+v", localsiteItem)
	}
	if prefix, ok := localsiteItem.Destination.(*frrProto.AccessListItem_IpPrefix); !ok {
		t.Errorf("Expected localsite item to be IP prefix, got: %+v", localsiteItem.Destination)
	} else {
		if prefix.IpPrefix.IpAddress != "192.168.1.0" || prefix.IpPrefix.PrefixLength != 24 {
			t.Errorf("Expected 192.168.1.0/24, got %s/%d", prefix.IpPrefix.IpAddress, prefix.IpPrefix.PrefixLength)
		}
	}
	// ^=========== ACCESS LIST TEST ==========^

	// v=========== ROUTE MAP TEST ==========v
	routeMap := config.RouteMap
	if len(routeMap) != 1 {
		t.Errorf("Expected 1 route map, got %d", len(routeMap))
	}
	rm, ok := routeMap["lanroutes"]
	if !ok {
		t.Errorf("Expected route map 'lanroutes' to exist")
	}
	if rm.Sequence != "10" {
		t.Errorf("Expected route map sequence '10', got '%s'", rm.Sequence)
	}
	if !rm.Permit {
		t.Error("Expected route map 'lanroutes' to be a permit rule")
	}
	if rm.Match != "ip address" || rm.AccessList != "localsite" {
		t.Errorf("Expected match 'ip address' with access list 'localsite', got match '%s' and access list '%s'",
			rm.Match, rm.AccessList)
	}
	// ^=========== ROUTE MAP TEST ==========^
}

//func TestParseConfig(t *testing.T) {
//	// Test valid config
//	tempDir := t.TempDir()
//	configPath := filepath.Join(tempDir, "ospf.conf")
//
//	configData := `! OSPF Configuration
//interface eth0
// ip ospf area 0.0.0.0
// ip ospf cost 10
//
//interface eth1
// ip ospf area 0.0.0.1
// ip ospf passive
//
//router ospf
// ospf router-id 192.168.1.1
// network 192.168.1.0/24 area 0.0.0.0
// network 192.168.2.0/24 area 0.0.0.1
//exit
//
//interface eth2
// ip ospf area 0.0.0.0
//`
//
//	err := os.WriteFile(configPath, []byte(configData), 0644)
//	if err != nil {
//		t.Fatalf("Failed to write test config file: %v", err)
//	}
//
//	config, err := aggregator.ParseStaticFRRConfig(configPath)
//	if err != nil {
//		t.Fatalf("ParseStaticFRRConfig failed: %v", err)
//	}
//
//	// Validate the parsed config
//	if config.RouterID != "192.168.1.1" {
//		t.Errorf("Expected RouterID to be '192.168.1.1', got '%s'", config.RouterID)
//	}
//
//	if len(config.Interfaces) != 3 {
//		t.Errorf("Expected 3 interfaces, got %d", len(config.Interfaces))
//	}
//
//	// Check eth0 interface
//	var eth0Found bool
//	for _, iface := range config.Interfaces {
//		if iface.Name == "eth0" {
//			eth0Found = true
//			if iface.Area != "0.0.0.0" {
//				t.Errorf("Expected eth0 area to be '0.0.0.0', got '%s'", iface.Area)
//			}
//			if iface.Cost != 10 {
//				t.Errorf("Expected eth0 cost to be 10, got %d", iface.Cost)
//			}
//			if iface.Passive {
//				t.Error("Expected eth0 to not be passive")
//			}
//		}
//	}
//	if !eth0Found {
//		t.Error("Interface eth0 not found in parsed config")
//	}
//
//	// Check eth1 interface
//	var eth1Found bool
//	for _, iface := range config.Interfaces {
//		if iface.Name == "eth1" {
//			eth1Found = true
//			if iface.Area != "0.0.0.1" {
//				t.Errorf("Expected eth1 area to be '0.0.0.1', got '%s'", iface.Area)
//			}
//			if !iface.Passive {
//				t.Error("Expected eth1 to be passive")
//			}
//		}
//	}
//	if !eth1Found {
//		t.Error("Interface eth1 not found in parsed config")
//	}
//
//	// Check areas
//	if len(config.Areas) != 2 {
//		t.Errorf("Expected 2 areas, got %d", len(config.Areas))
//	}
//
//	// Check area 0.0.0.0
//	var area0Found bool
//	for _, area := range config.Areas {
//		if area.ID == "0.0.0.0" {
//			area0Found = true
//			if len(area.Networks) != 1 {
//				t.Errorf("Expected 1 network in area 0.0.0.0, got %d", len(area.Networks))
//			}
//			if len(area.Networks) > 0 && area.Networks[0] != "192.168.1.0/24" {
//				t.Errorf("Expected network '192.168.1.0/24', got '%s'", area.Networks[0])
//			}
//		}
//	}
//	if !area0Found {
//		t.Error("Area 0.0.0.0 not found in parsed config")
//	}
//}

func TestParseConfigErrors(t *testing.T) {
	// Test nonexistent file
	_, err := aggregator.ParseStaticFRRConfig("nonexistent.conf")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestParseConfigEdgeCases(t *testing.T) {
	// Test empty config file
	tempDir := t.TempDir()
	emptyConfigPath := filepath.Join(tempDir, "empty.conf")

	err := os.WriteFile(emptyConfigPath, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to write empty config file: %v", err)
	}

	config, err := aggregator.ParseStaticFRRConfig(emptyConfigPath)
	if err != nil {
		t.Fatalf("ParseStaticFRRConfig failed for empty file: %v", err)
	}

	if config == nil {
		t.Fatal("Expected non-nil config for empty file")
	}

	if len(config.Interfaces) != 0 {
		t.Errorf("Expected 0 interfaces for empty file, got %d", len(config.Interfaces))
	}

	// Test config with only comments
	commentsConfigPath := filepath.Join(tempDir, "comments.conf")

	err = os.WriteFile(commentsConfigPath, []byte("! This is a comment\n! Another comment\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write comments config file: %v", err)
	}

	config, err = aggregator.ParseStaticFRRConfig(commentsConfigPath)
	if err != nil {
		t.Fatalf("ParseStaticFRRConfig failed for comments-only file: %v", err)
	}

	if config == nil {
		t.Fatal("Expected non-nil config for comments-only file")
	}

	// Test malformed interface line
	malformedConfigPath := filepath.Join(tempDir, "malformed.conf")

	err = os.WriteFile(malformedConfigPath, []byte("interface\n ip ospf area 0.0.0.0\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write malformed config file: %v", err)
	}

	config, err = aggregator.ParseStaticFRRConfig(malformedConfigPath)
	if err != nil {
		t.Fatalf("ParseStaticFRRConfig failed for malformed file: %v", err)
	}

	// The parser should be resilient to this kind of error
	if config == nil {
		t.Fatal("Expected non-nil config for malformed file")
	}
}
