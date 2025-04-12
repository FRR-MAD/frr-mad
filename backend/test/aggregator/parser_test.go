package aggregator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ba2025-ysmprc/frr-tui/backend/internal/aggregator"
)

func TestParseConfig(t *testing.T) {
	// Test valid config
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "ospf.conf")

	configData := `! OSPF Configuration
interface eth0
 ip ospf area 0.0.0.0
 ip ospf cost 10

interface eth1
 ip ospf area 0.0.0.1
 ip ospf passive

router ospf
 ospf router-id 192.168.1.1
 network 192.168.1.0/24 area 0.0.0.0
 network 192.168.2.0/24 area 0.0.0.1
exit

interface eth2
 ip ospf area 0.0.0.0
`

	err := os.WriteFile(configPath, []byte(configData), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	config, err := aggregator.ParseConfig(configPath)
	if err != nil {
		t.Fatalf("ParseConfig failed: %v", err)
	}

	// Validate the parsed config
	if config.RouterID != "192.168.1.1" {
		t.Errorf("Expected RouterID to be '192.168.1.1', got '%s'", config.RouterID)
	}

	if len(config.Interfaces) != 3 {
		t.Errorf("Expected 3 interfaces, got %d", len(config.Interfaces))
	}

	// Check eth0 interface
	var eth0Found bool
	for _, iface := range config.Interfaces {
		if iface.Name == "eth0" {
			eth0Found = true
			if iface.Area != "0.0.0.0" {
				t.Errorf("Expected eth0 area to be '0.0.0.0', got '%s'", iface.Area)
			}
			if iface.Cost != 10 {
				t.Errorf("Expected eth0 cost to be 10, got %d", iface.Cost)
			}
			if iface.Passive {
				t.Error("Expected eth0 to not be passive")
			}
		}
	}
	if !eth0Found {
		t.Error("Interface eth0 not found in parsed config")
	}

	// Check eth1 interface
	var eth1Found bool
	for _, iface := range config.Interfaces {
		if iface.Name == "eth1" {
			eth1Found = true
			if iface.Area != "0.0.0.1" {
				t.Errorf("Expected eth1 area to be '0.0.0.1', got '%s'", iface.Area)
			}
			if !iface.Passive {
				t.Error("Expected eth1 to be passive")
			}
		}
	}
	if !eth1Found {
		t.Error("Interface eth1 not found in parsed config")
	}

	// Check areas
	if len(config.Areas) != 2 {
		t.Errorf("Expected 2 areas, got %d", len(config.Areas))
	}

	// Check area 0.0.0.0
	var area0Found bool
	for _, area := range config.Areas {
		if area.ID == "0.0.0.0" {
			area0Found = true
			if len(area.Networks) != 1 {
				t.Errorf("Expected 1 network in area 0.0.0.0, got %d", len(area.Networks))
			}
			if len(area.Networks) > 0 && area.Networks[0] != "192.168.1.0/24" {
				t.Errorf("Expected network '192.168.1.0/24', got '%s'", area.Networks[0])
			}
		}
	}
	if !area0Found {
		t.Error("Area 0.0.0.0 not found in parsed config")
	}
}

func TestParseConfigErrors(t *testing.T) {
	// Test nonexistent file
	_, err := aggregator.ParseConfig("nonexistent.conf")
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

	config, err := aggregator.ParseConfig(emptyConfigPath)
	if err != nil {
		t.Fatalf("ParseConfig failed for empty file: %v", err)
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

	config, err = aggregator.ParseConfig(commentsConfigPath)
	if err != nil {
		t.Fatalf("ParseConfig failed for comments-only file: %v", err)
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

	config, err = aggregator.ParseConfig(malformedConfigPath)
	if err != nil {
		t.Fatalf("ParseConfig failed for malformed file: %v", err)
	}

	// The parser should be resilient to this kind of error
	if config == nil {
		t.Fatal("Expected non-nil config for malformed file")
	}
}
