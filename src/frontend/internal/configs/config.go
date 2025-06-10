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

type FrrMadTuiConfig struct {
	Pages map[string]PageConfig `mapstructure:"pages"`
}

type PageConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

type SocketConfig struct {
	UnixSocketLocation string `mapstructure:"unixsocketlocation"`
	UnixSocketName     string `mapstructure:"unixsocketname"`
	SocketType         string `mapstructure:"sockettype"`
}

type Config struct {
	Default   DefaultConfig   `mapstructure:"default"`
	FrrMadTui FrrMadTuiConfig `mapstructure:"frrmadtui"`
	Socket    SocketConfig    `mapstructure:"socket"`
}

func LoadConfig() (*Config, error) {
	tmpConf, ok := os.LookupEnv("FRR_MAD_CONFFILE")
	if ok {
		ConfigLocation = tmpConf
	}

	fmt.Println("Loading configuration file:", ConfigLocation)

	yamlPath := getYAMLPath()
	return loadYAMLConfig(yamlPath)
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
