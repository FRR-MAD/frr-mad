package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var ConfigLocation = "/etc/frr-mad/main.yaml"

type DefaultConfig struct {
	TempFiles  string `mapstructure:"tempfiles"`
	LogPath    string `mapstructure:"logpath"`
	ExportPath string `mapstructure:"exportpath"`
	DebugLevel string `mapstructure:"debuglevel"`
}

type SocketConfig struct {
	UnixSocketLocation string `mapstructure:"unixsocketlocation"`
	UnixSocketName     string `mapstructure:"unixsocketname"`
	SocketType         string `mapstructure:"sockettype"`
}

type AnalyzerConfig struct {
	Foo string `mapstructure:"foo"`
}

type AggregatorConfig struct {
	FRRMetricsURL   string `mapstructure:"frrmetricsurl"`
	FRRConfigPath   string `mapstructure:"frrconfigpath"`
	PollInterval    int    `mapstructure:"pollinterval"`
	SocketPathBGP   string `mapstructure:"socketpathbgp"`
	SocketPathOSPF  string `mapstructure:"socketpathospf"`
	SocketPathZebra string `mapstructure:"socketpathzebra"`
	SocketPath      string `mapstructure:"socketpath"`
}

type ExporterConfig struct {
	Foo string `mapstructure:"foo"`
}

type Config struct {
	Default    DefaultConfig    `mapstructure:"default"`
	Socket     SocketConfig     `mapstructure:"socket"`
	Analyzer   AnalyzerConfig   `mapstructure:"analyzer"`
	Aggregator AggregatorConfig `mapstructure:"aggregator"`
	Exporter   ExporterConfig   `mapstructure:"exporter"`
}

func LoadConfig() (*Config, error) {
	fmt.Println("Loading configuration file:", ConfigLocation)

	yamlPath := getYAMLPath()
	return loadYAMLConfig(yamlPath)
}

func GetConfig() map[string]map[string]string {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	result := make(map[string]map[string]string)

	result["default"] = map[string]string{
		"TempFiles":  config.Default.TempFiles,
		"LogPath":    config.Default.LogPath,
		"DebugLevel": config.Default.DebugLevel,
	}

	result["socket"] = map[string]string{
		"UnixSocketLocation": config.Socket.UnixSocketLocation,
		"UnixSocketName":     config.Socket.UnixSocketName,
		"SocketType":         config.Socket.SocketType,
	}

	result["analyzer"] = map[string]string{
		"foo": config.Analyzer.Foo,
	}

	result["aggregator"] = map[string]string{
		"FRRMetricsURL":   config.Aggregator.FRRMetricsURL,
		"FRRConfigPath":   config.Aggregator.FRRConfigPath,
		"PollInterval":    fmt.Sprintf("%d", config.Aggregator.PollInterval),
		"SocketPathBGP":   config.Aggregator.SocketPathBGP,
		"SocketPathOSPF":  config.Aggregator.SocketPathOSPF,
		"SocketPathZebra": config.Aggregator.SocketPathZebra,
		"SocketPath":      config.Aggregator.SocketPath,
	}

	result["exporter"] = map[string]string{
		"foo": config.Exporter.Foo,
	}

	return result
}

func loadYAMLConfig(yamlPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(yamlPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading YAML config: %w", err)
	}

	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return config, nil
}

func getYAMLPath() string {
	base := strings.TrimSuffix(ConfigLocation, filepath.Ext(ConfigLocation))
	return base + ".yaml"
}

func InitConfig() (*Config, error) {
	if _, err := os.Stat(ConfigLocation); os.IsNotExist(err) {
		yamlPath := getYAMLPath()
		if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
			return createDefaultConfig()
		}
	}

	return LoadConfig()
}

func createDefaultConfig() (*Config, error) {
	v := viper.New()

	configDir := filepath.Dir(getYAMLPath())
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	yamlPath := getYAMLPath()
	v.SetConfigFile(yamlPath)
	v.SetConfigType("yaml")

	if err := v.WriteConfig(); err != nil {
		return nil, fmt.Errorf("failed to write default config: %w", err)
	}

	fmt.Printf("Created default configuration at %s\n", yamlPath)

	return loadYAMLConfig(yamlPath)
}

func SaveConfig(config *Config) error {
	v := viper.New()

	v.Set("default.tempfiles", config.Default.TempFiles)
	v.Set("default.logpath", config.Default.LogPath)
	v.Set("default.debuglevel", config.Default.DebugLevel)
	v.Set("default.exportpath", config.Default.ExportPath)

	// Socket section
	v.Set("socket.unixsocketlocation", config.Socket.UnixSocketLocation)
	v.Set("socket.sockettype", config.Socket.SocketType)

	// Analyzer section
	v.Set("analyzer.foo", config.Analyzer.Foo)

	// Aggregator section
	v.Set("aggregator.frrmetricsurl", config.Aggregator.FRRMetricsURL)
	v.Set("aggregator.frrconfigpath", config.Aggregator.FRRConfigPath)
	v.Set("aggregator.pollinterval", config.Aggregator.PollInterval)
	v.Set("aggregator.socketpathbgp", config.Aggregator.SocketPathBGP)
	v.Set("aggregator.socketpathospf", config.Aggregator.SocketPathOSPF)
	v.Set("aggregator.socketpathzebra", config.Aggregator.SocketPathZebra)
	v.Set("aggregator.socketpath", config.Aggregator.SocketPath)

	// Exporter section
	v.Set("exporter.foo", config.Exporter.Foo)

	yamlPath := getYAMLPath()
	v.SetConfigFile(yamlPath)
	v.SetConfigType("yaml")

	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}
