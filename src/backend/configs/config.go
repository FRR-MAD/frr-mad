package configs

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var ConfigLocation = "/etc/frr-mad/frr-mad.yaml"

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
Port int                       `mapstructure:"Port"`             
OSPFRouterData string          `mapstructure:"OSPFRouterData"`              
OSPFAsbrSummaryData string     `mapstructure:"OSPFAsbrSummaryData"`                    
OSPFExternalData  string       `mapstructure:"OSPFExternalData"`
OSPFNssaExternalData  string   `mapstructure:"OSPFNssaExternalData"`
OSPFDatabase string            `mapstructure:"OSPFDatabase"`
OSPFDuplicates string          `mapstructure:"OSPFDuplicates"`
OSPFNeighbors  string          `mapstructure:"OSPFNeighbors"`
InterfaceList  string          `mapstructure:"OSPFInterfaceList"`
RouteList  string              `mapstructure:"OSPFRouteList"`
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

func getYAMLPath() string {
	base := strings.TrimSuffix(ConfigLocation, filepath.Ext(ConfigLocation))
	return base + ".yaml"
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


//func LoadConfig() map[string]map[string]string {
// func LoadConfig() (*Config, error) {
// 	fmt.Println("Loading configuration file:", ConfigLocation)

// 	yamlPath := getYAMLPath()
// 	return loadYAMLConfig(yamlPath)
// }
// func foobar() {}
// 	fmt.Println("Loading configuration file:", ConfigLocation)
// 	file, err := os.Open(ConfigLocation)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()

// 	config := make(map[string]map[string]string)
// 	scanner := bufio.NewScanner(file)
// 	var currentSection string

// 	for scanner.Scan() {
// 		line := strings.TrimSpace(scanner.Text())

// 		if line == "" || strings.HasPrefix(line, "#") {
// 			continue
// 		}

// 		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
// 			currentSection = line[1 : len(line)-1]
// 			config[currentSection] = make(map[string]string)
// 			continue
// 		}

// 		if currentSection != "" {
// 			parts := strings.SplitN(line, "=", 2)
// 			if len(parts) == 2 {
// 				key := strings.TrimSpace(parts[0])
// 				value := strings.TrimSpace(parts[1])
// 				config[currentSection][key] = value
// 			}
// 		}
// 	}

// 	if err := scanner.Err(); err != nil {
// 		panic(err)
// 	}

// 	return config
// }


//package configs

//import (
	//"bufio"
	//"fmt"
	//"os"
	//"path/filepath"
	//"strconv"
	//"strings"

	//"github.com/spf13/viper"
//)

//var ConfigLocation = "/etc/frr-analytics/main.conf"

//type DefaultConfig struct {
	//TempFiles  string `mapstructure:"tempfiles"`
	//LogPath    string `mapstructure:"logpath"`
	//DebugLevel string `mapstructure:"debuglevel"`
//}

//type SocketConfig struct {
	//UnixSocketLocation string `mapstructure:"unixsocketlocation"`
	//UnixSocketName     string `mapstructure:"unixsocketname"`
	//SocketType         string `mapstructure:"sockettype"`
//}

//type AnalyzerConfig struct {
	//Foo string `mapstructure:"foo"`
//}

//type AggregatorConfig struct {
	//FRRMetricsURL   string `mapstructure:"frrmetricsurl"`
	//FRRConfigPath   string `mapstructure:"frrconfigpath"`
	//PollInterval    int    `mapstructure:"pollinterval"`
	//SocketPathBGP   string `mapstructure:"socketpathbgp"`
	//SocketPathOSPF  string `mapstructure:"socketpathospf"`
	//SocketPathZebra string `mapstructure:"socketpathzebra"`
	//SocketPath      string `mapstructure:"socketpath"`
//}

//type ExporterConfig struct {
	//Foo string `mapstructure:"foo"`
//}

//type Config struct {
	//Default    DefaultConfig    `mapstructure:"default"`
	//Socket     SocketConfig     `mapstructure:"socket"`
	//Analyzer   AnalyzerConfig   `mapstructure:"analyzer"`
	//Aggregator AggregatorConfig `mapstructure:"aggregator"`
	//Exporter   ExporterConfig   `mapstructure:"exporter"`
//}

//func GetConfig() map[string]map[string]string {
	//config, _ := LoadConfig()

	////config = make(map[string]map[string]string)
	//scanner := bufio.NewScanner(file)
	//var currentSection string

	//for scanner.Scan() {
		//line := strings.TrimSpace(scanner.Text())

		//if line == "" || strings.HasPrefix(line, "#") {
			//continue
		//}

		//if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			//currentSection = line[1 : len(line)-1]
			//config[currentSection] = make(map[string]string)
			//continue
		//}

		//if currentSection != "" {
			//parts := strings.SplitN(line, "=", 2)
			//if len(parts) == 2 {
				//key := strings.TrimSpace(parts[0])
				//value := strings.TrimSpace(parts[1])
				//config[currentSection][key] = value
			//}
		//}
	//}

	//if err := scanner.Err(); err != nil {
		//panic(err)
	//}
//}

//type ParsedFlag struct {
	//Name        string
	//Description string
	//Enabled     bool
//}

//func LoadConfig() (*Config, error) {
	//fmt.Println("Loading configuration file:", ConfigLocation)


	//yamlPath := getYAMLPath()
	//return loadYAMLConfig(yamlPath)
//}

////func LoadConfig() map[string]map[string]string {
////	fmt.Println("Loading configuration file:", ConfigLocation)
////	file, err := os.Open(ConfigLocation)
////	if err != nil {
////		panic(err)
////	}
////	defer file.Close()
////
////	result := make(map[string]map[string]string)
////
////	result["default"] = map[string]string{
////		"TempFiles":  config.Default.TempFiles,
////		"LogPath":    config.Default.LogPath,
////		"DebugLevel": config.Default.DebugLevel,
////	}
////
////	result["socket"] = map[string]string{
////		"UnixSocketLocation": config.Socket.UnixSocketLocation,
////		"UnixSocketName":     config.Socket.UnixSocketName,
////		"SocketType":         config.Socket.SocketType,
////	}
////
////	result["analyzer"] = map[string]string{
////		"foo": config.Analyzer.Foo,
////	}
////
////	result["aggregator"] = map[string]string{
////		"FRRMetricsURL":   config.Aggregator.FRRMetricsURL,
////		"FRRConfigPath":   config.Aggregator.FRRConfigPath,
////		"PollInterval":    fmt.Sprintf("%d", config.Aggregator.PollInterval),
////		"SocketPathBGP":   config.Aggregator.SocketPathBGP,
////		"SocketPathOSPF":  config.Aggregator.SocketPathOSPF,
////		"SocketPathZebra": config.Aggregator.SocketPathZebra,
////		"SocketPath":      config.Aggregator.SocketPath,
////	}
////
////	result["exporter"] = map[string]string{
////		"foo": config.Exporter.Foo,
////	}
////
////	return result
////}

//func loadYAMLConfig(yamlPath string) (*Config, error) {
	//v := viper.New()
	//v.SetConfigFile(yamlPath)
	//v.SetConfigType("yaml")

	//if err := v.ReadInConfig(); err != nil {
		//return nil, fmt.Errorf("error reading YAML config: %w", err)
	//}

	//config := &Config{}
	//if err := v.Unmarshal(config); err != nil {
		//return nil, fmt.Errorf("error unmarshaling config: %w", err)
	//}

	//return config, nil
//}

//func getYAMLPath() string {
	//base := strings.TrimSuffix(ConfigLocation, filepath.Ext(ConfigLocation))
	//return base + ".yaml"
//}

//func InitConfig() (*Config, error) {
	//if _, err := os.Stat(ConfigLocation); os.IsNotExist(err) {
		//yamlPath := getYAMLPath()
		//if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
			//return createDefaultConfig()
		//}
	//}

	//return LoadConfig()
//}

//func createDefaultConfig() (*Config, error) {
	//v := viper.New()

	//configDir := filepath.Dir(getYAMLPath())
	//if err := os.MkdirAll(configDir, 0755); err != nil {
		//return nil, fmt.Errorf("failed to create config directory: %w", err)
	//}

	//yamlPath := getYAMLPath()
	//v.SetConfigFile(yamlPath)
	//v.SetConfigType("yaml")

	//if err := v.WriteConfig(); err != nil {
		//return nil, fmt.Errorf("failed to write default config: %w", err)
	//}

	//fmt.Printf("Created default configuration at %s\n", yamlPath)

	//return loadYAMLConfig(yamlPath)
//}

//func SaveConfig(config *Config) error {
	//v := viper.New()

	//v.Set("default.tempfiles", config.Default.TempFiles)
	//v.Set("default.logpath", config.Default.LogPath)
	//v.Set("default.debuglevel", config.Default.DebugLevel)

	//// Socket section
	//v.Set("socket.unixsocketlocation", config.Socket.UnixSocketLocation)
	//v.Set("socket.sockettype", config.Socket.SocketType)

	//// Analyzer section
	//v.Set("analyzer.foo", config.Analyzer.Foo)

	//// Aggregator section
	//v.Set("aggregator.frrmetricsurl", config.Aggregator.FRRMetricsURL)
	//v.Set("aggregator.frrconfigpath", config.Aggregator.FRRConfigPath)
	//v.Set("aggregator.pollinterval", config.Aggregator.PollInterval)
	//v.Set("aggregator.socketpathbgp", config.Aggregator.SocketPathBGP)
	//v.Set("aggregator.socketpathospf", config.Aggregator.SocketPathOSPF)
	//v.Set("aggregator.socketpathzebra", config.Aggregator.SocketPathZebra)
	//v.Set("aggregator.socketpath", config.Aggregator.SocketPath)

	//// Exporter section
	//v.Set("exporter.foo", config.Exporter.Foo)

	//yamlPath := getYAMLPath()
	//v.SetConfigFile(yamlPath)
	//v.SetConfigType("yaml")

	//if err := v.WriteConfig(); err != nil {
		//return fmt.Errorf("failed to save config: %w", err)
	//}

	//return nil
//}
//func ParseFlagTuple(tuple string) (*ParsedFlag, error) {
	//tuple = strings.TrimSpace(tuple)
	//if !strings.HasPrefix(tuple, "(") || !strings.HasSuffix(tuple, ")") {
		//return nil, fmt.Errorf("invalid flag tuple format - must be enclosed in parentheses")
	//}

	//tuple = tuple[1 : len(tuple)-1]

	//parts := splitTupleComponents(tuple)
	//if len(parts) != 3 {
		//return nil, fmt.Errorf("flag tuple must have exactly 3 components")
	//}

	//name := strings.TrimSpace(parts[0])
	//description := strings.TrimSpace(parts[1])
	//enabledStr := strings.TrimSpace(parts[2])

	//if strings.HasPrefix(description, `"`) && strings.HasSuffix(description, `"`) {
		//description = description[1 : len(description)-1]
	//}

	//enabled, err := strconv.ParseBool(enabledStr)
	//if err != nil {
		//return nil, fmt.Errorf("invalid boolean value in flag tuple: %v", err)
	//}

	//return &ParsedFlag{
		//Name:        name,
		//Description: description,
		//Enabled:     enabled,
	//}, nil
//}

//func splitTupleComponents(tuple string) []string {
	//var parts []string
	//var current strings.Builder
	//inQuotes := false

	//for _, r := range tuple {
		//switch {
		//case r == ',' && !inQuotes:
			//parts = append(parts, current.String())
			//current.Reset()
		//case r == '"':
			//inQuotes = !inQuotes
			//current.WriteRune(r)
		//default:
			//current.WriteRune(r)
		//}
	//}

	//if current.Len() > 0 {
		//parts = append(parts, current.String())
	//}

	//return parts
//}

//func GetFlagConfigs(config map[string]string) (map[string]ParsedFlag, error) {
	//flags := make(map[string]ParsedFlag)

	//for key, value := range config {
		//if strings.HasPrefix(value, "(") && strings.HasSuffix(value, ")") {
			//flag, err := ParseFlagTuple(value)
			//if err != nil {
				//return nil, fmt.Errorf("error parsing flag %s: %v", key, err)
			//}
			//flags[key] = *flag
		//}
	//}

	//return flags, nil
//}
