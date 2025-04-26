package aggregator_test

import (
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"os"
	"path/filepath"
	"testing"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/aggregator"
)

func TestParseStaticFRRConfig(t *testing.T) {
	configPathR101 := "./mock-files/r101.conf"
	configPathR103 := "./mock-files/r103.conf"
	configPathR112 := "./mock-files/r112.conf"
	configPathR203 := "./mock-files/r203.conf"
	configR101, err := aggregator.ParseStaticFRRConfig(configPathR101)
	if err != nil {
		t.Fatalf("ParseStaticFRRConfig configR101 failed: %v", err)
	}
	configR103, err := aggregator.ParseStaticFRRConfig(configPathR103)
	if err != nil {
		t.Fatalf("ParseStaticFRRConfig configR103 failed: %v", err)
	}
	configR112, err := aggregator.ParseStaticFRRConfig(configPathR112)
	if err != nil {
		t.Fatalf("ParseStaticFRRConfig configR112 failed: %v", err)
	}
	configR203, err := aggregator.ParseStaticFRRConfig(configPathR203)
	if err != nil {
		t.Fatalf("ParseStaticFRRConfig configR203 failed: %v", err)
	}

	t.Run("Metadata", func(t *testing.T) {
		testMetadataR101(t, configR101, configR103, configR112)
	})
	t.Run("Interfaces", func(t *testing.T) {
		testInterfaces(t, configR101, configR103, configR112, configR203)
	})
	t.Run("StaticRoutes", func(t *testing.T) {
		testStaticRoutes(t, configR101, configR103, configR112)
	})
	t.Run("OSPFConfig", func(t *testing.T) {
		testOSPFConfig(t, configR101, configR103, configR112)
	})
	t.Run("AccessList", func(t *testing.T) {
		testAccessList(t, configR101, configR103, configR112)
	})
	t.Run("RouteMap", func(t *testing.T) {
		testRouteMap(t, configR101, configR103, configR112)
	})
}

// vvvvvvvvvv HELPERS FOR TestParseStaticFRRConfig vvvvvvvvvv
func testMetadataR101(
	t *testing.T,
	configR101 *frrProto.StaticFRRConfiguration,
	configR103 *frrProto.StaticFRRConfiguration,
	configR112 *frrProto.StaticFRRConfiguration,
) {
	if configR101.Hostname != "r101" {
		t.Errorf("Expected Hostname to be 'r101', got '%s'", configR101.Hostname)
	}
	if configR101.FrrVersion != "8.5.4_git" {
		t.Errorf("Expected FrrVersion to be '8.5.4_git', got '%s'", configR101.FrrVersion)
	}
	if configR101.ServiceAdvancedVty != true {
		t.Error("Expected ServiceAdvancedVty to be true")
	}
	if configR103.Hostname != "r103" {
		t.Errorf("Expected Hostname to be 'r103', got '%s'", configR101.Hostname)
	}
	if configR112.Hostname != "r112" {
		t.Errorf("Expected Hostname to be 'r112', got '%s'", configR101.Hostname)
	}
}

func testInterfaces(
	t *testing.T,
	configR101 *frrProto.StaticFRRConfiguration,
	configR103 *frrProto.StaticFRRConfiguration,
	configR112 *frrProto.StaticFRRConfiguration,
	configR203 *frrProto.StaticFRRConfiguration,
) {
	// ========== r101 ==========
	if len(configR101.Interfaces) != 12 {
		t.Errorf("Expected 12 interfaces, got %d", len(configR101.Interfaces))
	}

	//r101Eth1
	r101Eth1 := configR101.Interfaces[0]
	if r101Eth1.Name != "eth1" {
		t.Errorf("Expected first interface name to be 'r101Eth1', got '%s'", r101Eth1.Name)
	}
	if len(r101Eth1.InterfaceIpPrefixes) == 0 {
		t.Error("Expected r101Eth1 to have at least one IP address")
	} else {
		if r101Eth1.InterfaceIpPrefixes[0].IpPrefix.IpAddress != "172.22.1.1" {
			t.Errorf("Expected r101Eth1 IP to be 172.22.1.1, got '%s'", r101Eth1.InterfaceIpPrefixes[0].IpPrefix.IpAddress)
		}
		if r101Eth1.InterfaceIpPrefixes[0].IpPrefix.PrefixLength != 24 {
			t.Errorf("Expected r101Eth1 prefix length to be 24, got %d", r101Eth1.InterfaceIpPrefixes[0].IpPrefix.PrefixLength)
		}
	}

	// r101Eth2
	r101Eth2 := configR101.Interfaces[1]
	if r101Eth2.Name != "eth2" {
		t.Errorf("Expected second interface name to be 'r101Eth2', got '%s'", r101Eth2.Name)
	}
	if len(r101Eth2.InterfaceIpPrefixes) != 2 {
		t.Errorf("Expected r101Eth2 to have 2 IP address, got '%v'", len(r101Eth2.InterfaceIpPrefixes))
	} else {
		if r101Eth2.InterfaceIpPrefixes[0].IpPrefix.IpAddress != "10.0.12.1" {
			t.Errorf("Expected r101Eth2 IP to be 10.0.12.1, got '%s'", r101Eth2.InterfaceIpPrefixes[0].IpPrefix.IpAddress)
		}
		if r101Eth2.InterfaceIpPrefixes[0].IpPrefix.PrefixLength != 24 {
			t.Errorf("Expected r101Eth2 prefix length to be 24, got %d", r101Eth2.InterfaceIpPrefixes[0].IpPrefix.PrefixLength)
		}
		if r101Eth2.InterfaceIpPrefixes[1].IpPrefix.IpAddress != "10.0.2.1" {
			t.Errorf("Expected r101Eth2 IP to be 10.0.2.1, got '%s'", r101Eth2.InterfaceIpPrefixes[1].IpPrefix.IpAddress)
		}
		if r101Eth2.InterfaceIpPrefixes[1].IpPrefix.PrefixLength != 24 {
			t.Errorf("Expected r101Eth2 prefix length to be 24, got %d", r101Eth2.InterfaceIpPrefixes[1].IpPrefix.PrefixLength)
		}
	}
	if r101Eth2.Area != "0.0.0.0" {
		t.Errorf("Expected r101Eth2 to be in Area '0.0.0.0', got '%s'", r101Eth2.Area)
	}
	if r101Eth2.InterfaceIpPrefixes[0].Passive {
		t.Errorf("Expected r101Eth2 IP Address to be in Passive 'false', got '%v'", r101Eth2.InterfaceIpPrefixes[0].Passive)
	}
	if !r101Eth2.InterfaceIpPrefixes[1].Passive {
		t.Errorf("Expected r101Eth2 IP Address '%s' to be in Passive 'false', got '%v'",
			r101Eth2.InterfaceIpPrefixes[1].IpPrefix.IpAddress,
			r101Eth2.InterfaceIpPrefixes[1].Passive)
	}

	// r101Eth3
	r101Eth3 := configR101.Interfaces[2]
	if r101Eth3.Name != "eth3" {
		t.Errorf("Expected r101Eth3 name to be 'r101Eth3', got '%s'", r101Eth3.Name)
	}
	if len(r101Eth3.InterfaceIpPrefixes) == 0 {
		t.Error("Expected r101Eth3 to have at least one IP address")
	} else {
		if r101Eth3.InterfaceIpPrefixes[0].IpPrefix.IpAddress != "10.0.13.1" {
			t.Errorf("Expected r101Eth3 IP to be 10.0.13.1, got '%s'", r101Eth3.InterfaceIpPrefixes[0].IpPrefix.IpAddress)
		}
		if r101Eth3.InterfaceIpPrefixes[0].IpPrefix.PrefixLength != 24 {
			t.Errorf("Expected r101Eth3 prefix length to be 24, got %d", r101Eth3.InterfaceIpPrefixes[0].IpPrefix.PrefixLength)
		}
	}
	if r101Eth3.Area != "0.0.0.0" {
		t.Errorf("Expected r101Eth3 to be in Area '0.0.0.0', got '%s'", r101Eth3.Area)
	}
	if r101Eth3.InterfaceIpPrefixes[0].Passive {
		t.Error("Expected r101Eth3 to be non-passive")
	}

	// r101Eth4
	r101Eth4 := configR101.Interfaces[3]
	if r101Eth4.Name != "eth4" {
		t.Errorf("Expected r101Eth4 name to be 'r101Eth4', got '%s'", r101Eth4.Name)
	}
	if len(r101Eth4.InterfaceIpPrefixes) == 0 {
		t.Error("Expected r101Eth4 to have at least one IP address")
	} else {
		if r101Eth4.InterfaceIpPrefixes[0].IpPrefix.IpAddress != "10.0.0.1" {
			t.Errorf("Expected r101Eth4 IP to be 10.0.0.1, got '%s'", r101Eth4.InterfaceIpPrefixes[0].IpPrefix.IpAddress)
		}
		if r101Eth4.InterfaceIpPrefixes[0].IpPrefix.PrefixLength != 23 {
			t.Errorf("Expected r101Eth4 prefix length to be 23, got %d", r101Eth4.InterfaceIpPrefixes[0].IpPrefix.PrefixLength)
		}
	}
	if r101Eth4.Area != "0.0.0.0" {
		t.Errorf("Expected r101Eth4 to be in Area '0.0.0.0', got '%s'", r101Eth4.Area)
	}
	if !r101Eth4.InterfaceIpPrefixes[0].Passive {
		t.Errorf("Expected r101Eth4 to be passive, got '%v'", r101Eth4.InterfaceIpPrefixes[0].Passive)
	}

	// r101Eth11
	r101Eth11 := configR101.Interfaces[10]
	if r101Eth11.Name != "eth11" {
		t.Errorf("Expected r101Eth11 name to be 'r101Eth11', got '%s'", r101Eth11.Name)
	}
	if len(r101Eth11.InterfaceIpPrefixes) == 0 {
		t.Error("Expected r101Eth11 to have at least one IP address")
	} else {
		if r101Eth11.InterfaceIpPrefixes[0].IpPrefix.IpAddress != "10.0.19.1" {
			t.Errorf("Expected r101Eth11 IP to be 10.0.19.1, got '%s'", r101Eth11.InterfaceIpPrefixes[0].IpPrefix.IpAddress)
		}
		if r101Eth11.InterfaceIpPrefixes[0].IpPrefix.PrefixLength != 24 {
			t.Errorf("Expected r101Eth11 prefix length to be 24, got %d", r101Eth11.InterfaceIpPrefixes[0].IpPrefix.PrefixLength)
		}
	}
	if r101Eth11.Area != "0.0.0.0" {
		t.Errorf("Expected r101Eth11 to be in Area '0.0.0.0', got '%s'", r101Eth11.Area)
	}
	if r101Eth11.InterfaceIpPrefixes[0].Passive {
		t.Error("Expected r101Eth11 to be non-passive")
	}

	// lo
	r101Lo := configR101.Interfaces[11]
	if r101Lo.Name != "lo" {
		t.Errorf("Expected r101Lo interface name to be 'r101Lo', got '%s'", r101Lo.Name)
	}
	if len(r101Lo.InterfaceIpPrefixes) == 0 {
		t.Error("Expected r101Lo to have at least one IP address")
	} else {
		if r101Lo.InterfaceIpPrefixes[0].IpPrefix.IpAddress != "65.0.1.1" {
			t.Errorf("Expected r101Lo IP to be 65.0.1.1, got '%s'", r101Lo.InterfaceIpPrefixes[0].IpPrefix.IpAddress)
		}
		if r101Lo.InterfaceIpPrefixes[0].IpPrefix.PrefixLength != 32 {
			t.Errorf("Expected r101Lo prefix length to be 32, got %d", r101Lo.InterfaceIpPrefixes[0].IpPrefix.PrefixLength)
		}
	}
	if !r101Lo.InterfaceIpPrefixes[0].Passive {
		t.Error("Expected r101Lo to be passive")
	}

	// ========== r103 ==========
	if len(configR103.Interfaces) != 4 {
		t.Errorf("Expected 4 interfaces, got %d", len(configR103.Interfaces))
	}

	// r103Eth1
	r103Eth1 := configR103.Interfaces[0]
	if r103Eth1.Name != "eth1" {
		t.Errorf("Expected first interface name to be 'r103Eth1', got '%s'", r103Eth1.Name)
	}
	if len(r103Eth1.InterfaceIpPrefixes) != 3 {
		t.Errorf("Expected r103Eth1 to have 3 IP address, got '%v'", len(r103Eth1.InterfaceIpPrefixes))
	} else {
		if r103Eth1.InterfaceIpPrefixes[0].IpPrefix.IpAddress != "10.0.13.3" {
			t.Errorf("Expected r103Eth1 IP1 to be 10.0.13.3, got '%s'", r103Eth1.InterfaceIpPrefixes[0].IpPrefix.IpAddress)
		}
		if r103Eth1.InterfaceIpPrefixes[0].IpPrefix.PrefixLength != 24 {
			t.Errorf("Expected r103Eth1 prefix length to be 24, got %d", r103Eth1.InterfaceIpPrefixes[0].IpPrefix.PrefixLength)
		}
		if r103Eth1.InterfaceIpPrefixes[1].IpPrefix.IpAddress != "10.0.13.33" {
			t.Errorf("Expected r103Eth1 IP2 to be 10.0.13.33, got '%s'", r103Eth1.InterfaceIpPrefixes[1].IpPrefix.IpAddress)
		}
		if r103Eth1.InterfaceIpPrefixes[2].IpPrefix.IpAddress != "10.0.13.30" {
			t.Errorf("Expected r103Eth1 IP3 to be 10.0.13.30, got '%s'", r103Eth1.InterfaceIpPrefixes[2].IpPrefix.IpAddress)
		}
	}
	if r103Eth1.Area != "0.0.0.0" {
		t.Errorf("Expected r103Eth1 to be in Area '0.0.0.0', got '%s'", r103Eth1.Area)
	}
	if r103Eth1.InterfaceIpPrefixes[0].Passive {
		t.Errorf("Expected r103Eth1 IP Address '%s' to be non-passive", r103Eth1.InterfaceIpPrefixes[0].IpPrefix.IpAddress)
	}
	if r103Eth1.InterfaceIpPrefixes[1].Passive {
		t.Errorf("Expected r103Eth1 IP Address '%s' to be non-passive", r103Eth1.InterfaceIpPrefixes[1].IpPrefix.IpAddress)
	}
	if r103Eth1.InterfaceIpPrefixes[1].Passive {
		t.Errorf("Expected r103Eth1 IP Address '%s' to be non-passive", r103Eth1.InterfaceIpPrefixes[2].IpPrefix.IpAddress)
	}

	// r103Eth2
	r103Eth2 := configR103.Interfaces[1]
	if r103Eth2.Name != "eth2" {
		t.Errorf("Expected second interface name to be 'eth2', got '%s'", r103Eth2.Name)
	}
	if len(r103Eth2.InterfaceIpPrefixes) == 0 {
		t.Error("Expected r103Eth2 to have at least one IP address")
	} else {
		if r103Eth2.InterfaceIpPrefixes[0].IpPrefix.IpAddress != "10.0.23.3" {
			t.Errorf("Expected r103Eth2 IP to be 10.0.23.3, got '%s'", r103Eth2.InterfaceIpPrefixes[0].IpPrefix.IpAddress)
		}
		if r103Eth2.InterfaceIpPrefixes[0].IpPrefix.PrefixLength != 24 {
			t.Errorf("Expected r103Eth2 prefix length to be 24, got %d", r103Eth2.InterfaceIpPrefixes[0].IpPrefix.PrefixLength)
		}
	}
	if r103Eth2.Area != "0.0.0.0" {
		t.Errorf("Expected r103Eth2 to be in Area '0.0.0.0', got '%s'", r103Eth2.Area)
	}
	if r103Eth2.InterfaceIpPrefixes[0].Passive {
		t.Error("Expected r103Eth2 to be non-passive")
	}

	// r103Eth3
	r103Eth3 := configR103.Interfaces[2]
	if r103Eth3.Name != "eth3" {
		t.Errorf("Expected second interface name to be 'eth3', got '%s'", r103Eth3.Name)
	}
	if len(r103Eth3.InterfaceIpPrefixes) == 0 {
		t.Error("Expected r103Eth3 to have at least one IP address")
	} else {
		if r103Eth3.InterfaceIpPrefixes[0].IpPrefix.IpAddress != "10.2.31.3" {
			t.Errorf("Expected r103Eth3 IP to be 10.2.31.3, got '%s'", r103Eth3.InterfaceIpPrefixes[0].IpPrefix.IpAddress)
		}
		if r103Eth3.InterfaceIpPrefixes[0].IpPrefix.PrefixLength != 24 {
			t.Errorf("Expected r103Eth3 prefix length to be 24, got %d", r103Eth3.InterfaceIpPrefixes[0].IpPrefix.PrefixLength)
		}
	}
	if r103Eth3.Area != "0.0.0.2" {
		t.Errorf("Expected r103Eth3 to be in Area '0.0.0.0', got '%s'", r103Eth3.Area)
	}
	if r103Eth3.InterfaceIpPrefixes[0].Passive {
		t.Error("Expected r103Eth3 to be non-passive")
	}

	// ========== r112 ==========
	if len(configR112.Interfaces) != 3 {
		t.Errorf("Expected 3 interfaces, got %d", len(configR112.Interfaces))
	}

	// r112Eth1
	r112Eth1 := configR112.Interfaces[0]
	if r112Eth1.Name != "eth1" {
		t.Errorf("Expected first interface name to be 'eth1', got '%s'", r112Eth1.Name)
	}
	if len(r112Eth1.InterfaceIpPrefixes) != 1 {
		t.Errorf("Expected r112Eth1 to have 1 IP address, got '%v'", len(r112Eth1.InterfaceIpPrefixes))
	} else {
		if r112Eth1.InterfaceIpPrefixes[0].IpPrefix.IpAddress != "10.1.12.12" {
			t.Errorf("Expected r112Eth1 IP1 to be 10.1.12.12, got '%s'", r112Eth1.InterfaceIpPrefixes[0].IpPrefix.IpAddress)
		}
		if r112Eth1.InterfaceIpPrefixes[0].IpPrefix.PrefixLength != 24 {
			t.Errorf("Expected r112Eth1 prefix length to be 24, got %d", r112Eth1.InterfaceIpPrefixes[0].IpPrefix.PrefixLength)
		}
	}
	if r112Eth1.Area != "0.0.0.1" {
		t.Errorf("Expected r112Eth1 to be in Area '0.0.0.1', got '%s'", r112Eth1.Area)
	}
	if r112Eth1.InterfaceIpPrefixes[0].Passive {
		t.Error("Expected r112Eth1 to be non-passive")
	}

	// ========== r203 ==========
	r203Eth1 := configR203.Interfaces[0]
	if r203Eth1.Name != "eth1" {
		t.Errorf("Expected first interface name to be 'eth1', got '%s'", r203Eth1.Name)
	}
	if r203Eth1.InterfaceIpPrefixes[0].IpPrefix.IpAddress != "10.20.13.3" {
		t.Errorf("Expected r203Eth1 IP1 to be 10.20.13.3, got '%s'", r203Eth1.InterfaceIpPrefixes[0].IpPrefix.IpAddress)
	}
	if r203Eth1.InterfaceIpPrefixes[0].IpPrefix.PrefixLength != 32 {
		t.Errorf("Expected r203Eth1 prefix length to be 32, got %d", r203Eth1.InterfaceIpPrefixes[0].IpPrefix.PrefixLength)
	}
	if r203Eth1.InterfaceIpPrefixes[0].PeerIpPrefix.IpAddress != "10.20.13.1" {
		t.Errorf("Expected r203Eth1 peer IP to be 10.20.13.1, got '%s'",
			r203Eth1.InterfaceIpPrefixes[0].PeerIpPrefix.IpAddress)
	}
	if r203Eth1.InterfaceIpPrefixes[0].PeerIpPrefix.PrefixLength != 32 {
		t.Errorf("Expected r203Eth1 peer IP prefix length to be 32, got %d",
			r203Eth1.InterfaceIpPrefixes[0].PeerIpPrefix.PrefixLength)
	}
}

func testStaticRoutes(
	t *testing.T,
	configR101 *frrProto.StaticFRRConfiguration,
	configR103 *frrProto.StaticFRRConfiguration,
	configR112 *frrProto.StaticFRRConfiguration,
) {
	// ========== r101 ==========
	staticRoutesR101 := configR101.StaticRoutes
	if len(staticRoutesR101) != 3 {
		t.Errorf("Expected 1 static route, got %d", len(staticRoutesR101))
	}
	// static route one
	if staticRoutesR101[0].IpPrefix.IpAddress != "192.168.1.0" {
		t.Errorf("Expected static route one to have 192.168.1.0, got '%s'", staticRoutesR101[0].IpPrefix)
	}
	if staticRoutesR101[0].IpPrefix.PrefixLength != 24 {
		t.Errorf("Expected static route Prefix one to be 24, got %d", staticRoutesR101[0].IpPrefix.PrefixLength)
	}
	if staticRoutesR101[0].NextHop != "192.168.100.91" {
		t.Errorf("Expected Next Hop to be 192.168.100.91, got '%s'", staticRoutesR101[0].NextHop)
	}
	// static route two
	if staticRoutesR101[1].IpPrefix.IpAddress != "192.168.2.0" {
		t.Errorf("Expected static route one to have 192.168.2.0, got '%s'", staticRoutesR101[1].IpPrefix)
	}
	if staticRoutesR101[1].IpPrefix.PrefixLength != 23 {
		t.Errorf("Expected static route Prefix one to be 23, got %d", staticRoutesR101[1].IpPrefix.PrefixLength)
	}
	if staticRoutesR101[1].NextHop != "192.168.102.91" {
		t.Errorf("Expected Next Hop to be 192.168.102.91, got '%s'", staticRoutesR101[1].NextHop)
	}
	// static route three
	if staticRoutesR101[2].IpPrefix.IpAddress != "192.168.4.0" {
		t.Errorf("Expected static route one to have 192.168.4.0, got '%s'", staticRoutesR101[2].IpPrefix)
	}
	if staticRoutesR101[2].IpPrefix.PrefixLength != 22 {
		t.Errorf("Expected static route Prefix one to be 22, got %d", staticRoutesR101[2].IpPrefix.PrefixLength)
	}
	if staticRoutesR101[2].NextHop != "192.168.104.91" {
		t.Errorf("Expected Next Hop to be 192.168.104.91, got '%s'", staticRoutesR101[2].NextHop)
	}
}

func testOSPFConfig(
	t *testing.T,
	configR101 *frrProto.StaticFRRConfiguration,
	configR103 *frrProto.StaticFRRConfiguration,
	configR112 *frrProto.StaticFRRConfiguration,
) {
	// ========== r101 ==========
	ospfConfigR101 := configR101.OspfConfig
	if ospfConfigR101.RouterId != "65.0.1.1" {
		t.Errorf("Expected ospf router id 65.0.1.1, got '%s'", ospfConfigR101.RouterId)
	}
	if ospfConfigR101.Redistribution[0].Type != "static" {
		t.Errorf("Expected ospf redistribution type static, got '%s'", ospfConfigR101.Redistribution[0].Type)
	}
	if ospfConfigR101.Redistribution[0].Metric != "1" {
		t.Errorf("Expected Metric to be '1', got '%s'", ospfConfigR101.Redistribution[0].Metric)
	}
	if ospfConfigR101.Redistribution[0].RouteMap != "lanroutes" {
		t.Errorf("Expected RouteMap to be 'lanroutes', got '%s'", ospfConfigR101.Redistribution[0].RouteMap)
	}
	if ospfConfigR101.Redistribution[1].Type != "bgp" {
		t.Errorf("Expected Redistribution type bgp, got '%s'", ospfConfigR101.Redistribution[1].Type)
	}
	if ospfConfigR101.Redistribution[1].Metric != "1" {
		t.Errorf("Expected Metric to be '1', got '%s'", ospfConfigR101.Redistribution[1].Metric)
	}
	if ospfConfigR101.Redistribution[1].RouteMap != "" {
		t.Errorf("Expected RouteMap to be empty, got '%s'", ospfConfigR101.Redistribution[1].RouteMap)
	}
	// ========== r103 ==========
	ospfConfigR103 := configR103.OspfConfig
	if ospfConfigR103.RouterId != "65.0.1.3" {
		t.Errorf("Expected ospf router id 65.0.1.3, got '%s'", ospfConfigR103.RouterId)
	}
	if ospfConfigR103.Area[0].Name != "0.0.0.2" {
		t.Errorf("Expected ospf area type transit, got '%s'", ospfConfigR103.Area[0].Type)
	}
	if ospfConfigR103.Area[0].Type != "transit" {
		t.Errorf("Expected ospf area type transit, got '%s'", ospfConfigR103.Area[0].Type)
	}
	if ospfConfigR103.VirtualLinkNeighbor != "65.0.1.22" {
		t.Errorf("Expected ospf virtual link neighbor 65.0.1.22, got '%s'", ospfConfigR103.VirtualLinkNeighbor)
	}

	// ========== r112 ==========
	ospfConfigR112 := configR112.OspfConfig
	if ospfConfigR112.RouterId != "65.0.1.12" {
		t.Errorf("Expected ospf router id 65.0.1.12, got '%s'", ospfConfigR112.RouterId)
	}
	if ospfConfigR112.Area[0].Name != "0.0.0.1" {
		t.Errorf("Expected ospf area type transit, got '%s'", ospfConfigR112.Area[0].Type)
	}
	if ospfConfigR112.Area[0].Type != "nssa" {
		t.Errorf("Expected ospf area type nssa, got '%s'", ospfConfigR112.Area[0].Type)
	}
	if ospfConfigR112.Redistribution[0].Type != "bgp" {
		t.Errorf("Expected Redistribution type bgp, got '%s'", ospfConfigR112.Redistribution[0].Type)
	}
	if ospfConfigR112.Redistribution[0].Metric != "1" {
		t.Errorf("Expected Metric to be '1', got '%s'", ospfConfigR112.Redistribution[0].Metric)
	}
	if ospfConfigR112.Redistribution[0].RouteMap != "" {
		t.Errorf("Expected RouteMap to be empty, got '%s'", ospfConfigR112.Redistribution[0].RouteMap)
	}
}

func testAccessList(
	t *testing.T,
	configR101 *frrProto.StaticFRRConfiguration,
	configR103 *frrProto.StaticFRRConfiguration,
	configR112 *frrProto.StaticFRRConfiguration,
) {
	// ========== r101 ==========
	accessList := configR101.AccessList
	if len(accessList) != 2 {
		t.Errorf("Expected 2 access list, got %d", len(accessList))
	}
	termList, ok := configR101.AccessList["term"]
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

	localsiteList, ok := configR101.AccessList["localsite"]
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
}

func testRouteMap(
	t *testing.T,
	configR101 *frrProto.StaticFRRConfiguration,
	configR103 *frrProto.StaticFRRConfiguration,
	configR112 *frrProto.StaticFRRConfiguration,
) {
	routeMap := configR101.RouteMap
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
}

// ^^^^^^^^^^ HELPERS FOR TestParseStaticFRRConfig ^^^^^^^^^^

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
