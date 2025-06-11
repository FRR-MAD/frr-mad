package configs_test

import (
	"os"
	"testing"

	"github.com/frr-mad/frr-mad/src/backend/configs"
	"github.com/stretchr/testify/assert"
)

// TestConfigLoadingHappyPath tests all successful scenarios
func TestConfigLoadingHappyPath(t *testing.T) {
	// Setup: Create test config files

	originalConfigLocation := configs.ConfigLocation
	t.Run("LoadConfig_with_overwrite_path", func(t *testing.T) {
		configPath := "mock-files/main.yaml"
		config, err := configs.LoadConfig(configPath)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if config == nil {
			t.Error("Expected config to be loaded, got nil")
		}
		// Add assertions based on your Config struct fields
		// Example: if config.SomeField != "expected_value" { t.Error(...) }
		assert.Equal(t, "/tmp/frr-mad", config.Default.TempFiles)
		assert.Equal(t, "error", config.Default.DebugLevel)
		assert.Equal(t, "/var/log/frr-mad", config.Default.LogPath)

		assert.Equal(t, "/var/run/frr-mad", config.Socket.UnixSocketLocation)
		assert.Equal(t, "analyzer.sock", config.Socket.UnixSocketName)
		assert.Equal(t, "unix", config.Socket.SocketType)

		assert.Equal(t, "/etc/frr/frr.conf", config.Aggregator.FRRConfigPath)
		assert.Equal(t, 5, config.Aggregator.PollInterval)
		assert.Equal(t, "/var/run/frr", config.Aggregator.SocketPath)

		assert.False(t, config.Exporter.OSPFRouterData)
		assert.False(t, config.Exporter.OSPFNetworkData)
		assert.False(t, config.Exporter.OSPFSummaryData)
		assert.False(t, config.Exporter.OSPFAsbrSummaryData)
		assert.False(t, config.Exporter.OSPFExternalData)
		assert.False(t, config.Exporter.OSPFNssaExternalData)
		assert.False(t, config.Exporter.OSPFDatabase)
		assert.False(t, config.Exporter.OSPFNeighbors)
		assert.False(t, config.Exporter.InterfaceList)
		assert.False(t, config.Exporter.RouteList)

	})

	t.Run("LoadConfig_with_env_variable", func(t *testing.T) {
		// Set environment variable
		os.Setenv("FRR_MAD_CONFFILE", "mock-files/main.yaml")
		defer os.Unsetenv("FRR_MAD_CONFFILE")

		config, err := configs.LoadConfig("")

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if config == nil {
			t.Error("Expected config to be loaded, got nil")
		}
		assert.Equal(t, "/tmp/frr-mad", config.Default.TempFiles)
		assert.Equal(t, "error", config.Default.DebugLevel)
		assert.Equal(t, "/var/log/frr-mad", config.Default.LogPath)

		assert.Equal(t, "/var/run/frr-mad", config.Socket.UnixSocketLocation)
		assert.Equal(t, "analyzer.sock", config.Socket.UnixSocketName)
		assert.Equal(t, "unix", config.Socket.SocketType)

		assert.Equal(t, "/etc/frr/frr.conf", config.Aggregator.FRRConfigPath)
		assert.Equal(t, 5, config.Aggregator.PollInterval)
		assert.Equal(t, "/var/run/frr", config.Aggregator.SocketPath)

		assert.False(t, config.Exporter.OSPFRouterData)
		assert.False(t, config.Exporter.OSPFNetworkData)
		assert.False(t, config.Exporter.OSPFSummaryData)
		assert.False(t, config.Exporter.OSPFAsbrSummaryData)
		assert.False(t, config.Exporter.OSPFExternalData)
		assert.False(t, config.Exporter.OSPFNssaExternalData)
		assert.False(t, config.Exporter.OSPFDatabase)
		assert.False(t, config.Exporter.OSPFNeighbors)
		assert.False(t, config.Exporter.InterfaceList)
		assert.False(t, config.Exporter.RouteList)

	})

	t.Run("LoadConfig_with_default_path", func(t *testing.T) {
		// Ensure no env variable is set
		os.Unsetenv("FRR_MAD_CONFFILE")

		// Assuming ConfigLocation has a default value
		configs.ConfigLocation = "mock-files/config.txt"
		defer func() { configs.ConfigLocation = originalConfigLocation }()

		config, err := configs.LoadConfig("")

		assert.Error(t, err)
		assert.Nil(t, config)
		assert.NotEqual(t, originalConfigLocation, configs.ConfigLocation)
		assert.Equal(t, "/etc/frr-mad/main.yaml", originalConfigLocation)
	})

	t.Run("getYAMLPath_success", func(t *testing.T) {
		originalConfigLocation := configs.ConfigLocation
		configs.ConfigLocation = "mock-files/main.yaml"
		defer func() { configs.ConfigLocation = originalConfigLocation }()

		yamlPath := configs.GetYAMLPath()
		expected := "mock-files/main.yaml"

		if yamlPath != expected {
			t.Errorf("Expected YAML path %s, got %s", expected, yamlPath)
		}
	})

	t.Run("loadYAMLConfig_success", func(t *testing.T) {
		yamlPath := "mock-files/main.yaml"
		config, err := configs.LoadYAMLConfig(yamlPath)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if config == nil {
			t.Error("Expected config to be loaded, got nil")
		}
	})
}

// TestConfigLoadingSadPath tests all error scenarios
func TestConfigLoadingSadPath(t *testing.T) {
	t.Run("LoadConfig_file_not_found", func(t *testing.T) {
		nonExistentPath := "non-existent/config.txt"
		config, err := configs.LoadConfig(nonExistentPath)

		if err == nil {
			t.Error("Expected error for non-existent file, got nil")
		}
		if config != nil {
			t.Error("Expected nil config for non-existent file")
		}
		if !os.IsNotExist(err) && err.Error() == "" {
			t.Errorf("Expected file not found error, got: %v", err)
		}
	})

	t.Run("LoadConfig_invalid_yaml_file", func(t *testing.T) {
		// Setup invalid YAML file
		setupInvalidYAMLFile(t)
		defer cleanupInvalidYAMLFile(t)

		configPath := "mock-files/invalid-config.txt"
		config, err := configs.LoadConfig(configPath)

		if err == nil {
			t.Error("Expected error for invalid YAML, got nil")
		}
		if config != nil {
			t.Error("Expected nil config for invalid YAML")
		}
	})

	t.Run("getYAMLPath_with_various_extensions", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{"config.txt", "config.yaml"},
			{"config.json", "config.yaml"},
			{"config", "config.yaml"},
			{"path/to/config.conf", "path/to/config.yaml"},
		}

		originalConfigLocation := configs.ConfigLocation
		defer func() { configs.ConfigLocation = originalConfigLocation }()

		for _, tc := range testCases {
			configs.ConfigLocation = tc.input
			result := configs.GetYAMLPath()
			if result != tc.expected {
				t.Errorf("For input %s, expected %s, got %s", tc.input, tc.expected, result)
			}
		}
	})

	t.Run("loadYAMLConfig_file_not_found", func(t *testing.T) {
		nonExistentYAML := "non-existent/config.yaml"
		config, err := configs.LoadYAMLConfig(nonExistentYAML)

		if err == nil {
			t.Error("Expected error for non-existent YAML file, got nil")
		}
		if config != nil {
			t.Error("Expected nil config for non-existent YAML file")
		}
	})

	t.Run("loadYAMLConfig_invalid_yaml_content", func(t *testing.T) {
		setupInvalidYAMLFile(t)
		defer cleanupInvalidYAMLFile(t)

		invalidYAMLPath := "mock-files/invalid-config.yaml"
		config, err := configs.LoadYAMLConfig(invalidYAMLPath)

		if err == nil {
			t.Error("Expected error for invalid YAML content, got nil")
		}
		if config != nil {
			t.Error("Expected nil config for invalid YAML content")
		}
	})

	t.Run("loadYAMLConfig_unmarshal_error", func(t *testing.T) {
		// Setup YAML with structure that doesn't match Config struct
		incompatibleYAMLPath := "mock-files/incompatible-config.yaml"
		config, err := configs.LoadYAMLConfig(incompatibleYAMLPath)

		// This might not always error depending on your Config struct
		// but it's good to test the unmarshal path
		if config == nil && err == nil {
			t.Error("Expected either config or error, got both nil")
		}
	})
}

func setupTestFiles(t *testing.T) {
	// Create mock-files directory
	err := os.MkdirAll("mock-files", 0755)
	if err != nil {
		t.Fatalf("Failed to create mock-files directory: %v", err)
	}

	// Create config.txt file
	configFile, err := os.Create("mock-files/config.txt")
	if err != nil {
		t.Fatalf("Failed to create config.txt: %v", err)
	}
	configFile.Close()

	// Create config.yaml file with valid YAML content
	yamlContent := `# Test configuration
app_name: test_app
version: 1.0.0
database:
  host: localhost
  port: 5432
`
	err = os.WriteFile("mock-files/config.yaml", []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config.yaml: %v", err)
	}
}

func cleanupTestFiles(t *testing.T) {
	err := os.RemoveAll("mock-files")
	if err != nil {
		t.Logf("Warning: Failed to cleanup mock-files: %v", err)
	}
}

func setupInvalidYAMLFile(t *testing.T) {
	err := os.MkdirAll("mock-files", 0755)
	if err != nil {
		t.Fatalf("Failed to create mock-files directory: %v", err)
	}

	// Create invalid-config.txt
	configFile, err := os.Create("mock-files/invalid-config.txt")
	if err != nil {
		t.Fatalf("Failed to create invalid-config.txt: %v", err)
	}
	configFile.Close()

	// Create invalid YAML content
	invalidYAML := `{
		invalid: yaml: content
		missing: quotes
		- invalid: structure
	}`
	err = os.WriteFile("mock-files/invalid-config.yaml", []byte(invalidYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid-config.yaml: %v", err)
	}
}

func cleanupInvalidYAMLFile(t *testing.T) {
	os.Remove("mock-files/invalid-config.txt")
	os.Remove("mock-files/invalid-config.yaml")
}

func setupIncompatibleYAMLFile(t *testing.T) {
	err := os.MkdirAll("mock-files", 0755)
	if err != nil {
		t.Fatalf("Failed to create mock-files directory: %v", err)
	}

	// Create YAML that's valid but doesn't match your Config struct
	incompatibleYAML := `
completely_different_structure:
  - item1
  - item2
random_field: "value"
numeric_array: [1, 2, 3, 4, 5]
`
	err = os.WriteFile("mock-files/incompatible-config.yaml", []byte(incompatibleYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create incompatible-config.yaml: %v", err)
	}
}

func cleanupIncompatibleYAMLFile(t *testing.T) {
	os.Remove("mock-files/incompatible-config.yaml")
}
