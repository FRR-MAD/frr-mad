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
	DebugLevel string `mapstructure:"debuglevel"`
}

type SocketConfig struct {
	UnixSocketLocation string `mapstructure:"unixsocketlocation"`
	UnixSocketName     string `mapstructure:"unixsocketname"`
	SocketType         string `mapstructure:"sockettype"`
}

type AggregatorConfig struct {
	FRRConfigPath string `mapstructure:"frrconfigpath"`
	PollInterval  int    `mapstructure:"pollinterval"`
	SocketPath    string `mapstructure:"socketpath"`
}

type ExporterConfig struct {
	Port                 int  `mapstructure:"Port"`
	OSPFRouterData       bool `mapstructure:"OSPFRouterData"`
	OSPFNetworkData      bool `mapstructure:"OSPFNetworkData"`
	OSPFSummaryData      bool `mapstructure:"OSPFSummaryData"`
	OSPFAsbrSummaryData  bool `mapstructure:"OSPFAsbrSummaryData"`
	OSPFExternalData     bool `mapstructure:"OSPFExternalData"`
	OSPFNssaExternalData bool `mapstructure:"OSPFNssaExternalData"`
	OSPFDatabase         bool `mapstructure:"OSPFDatabase"`
	OSPFNeighbors        bool `mapstructure:"OSPFNeighbors"`
	InterfaceList        bool `mapstructure:"InterfaceList"`
	RouteList            bool `mapstructure:"RouteList"`
}

type Config struct {
	Default    DefaultConfig    `mapstructure:"default"`
	Socket     SocketConfig     `mapstructure:"socket"`
	Aggregator AggregatorConfig `mapstructure:"aggregator"`
	Exporter   ExporterConfig   `mapstructure:"exporter"`
}

func LoadConfig(overwriteConfigPath string) (*Config, error) {
	if overwriteConfigPath != "" {
		ConfigLocation = overwriteConfigPath
	}

	tmpConf, ok := os.LookupEnv("FRR_MAD_CONFFILE")
	if ok {
		ConfigLocation = tmpConf
	}

	file, err := os.Open(ConfigLocation)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	yamlPath := GetYAMLPath()
	result, err := LoadYAMLConfig(yamlPath)
	return result, err
}

func GetYAMLPath() string {

	base := strings.TrimSuffix(ConfigLocation, filepath.Ext(ConfigLocation))
	return base + ".yaml"
}

func LoadYAMLConfig(yamlPath string) (*Config, error) {
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
