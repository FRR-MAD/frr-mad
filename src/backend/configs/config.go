package configs

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var ConfigLocation = "/etc/frr-analytics/main.conf"

type ParsedFlag struct {
	Name        string
	Description string
	Enabled     bool
}

func LoadConfig() map[string]map[string]string {
	fmt.Println("Loading configuration file:", ConfigLocation)
	file, err := os.Open(ConfigLocation)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	config := make(map[string]map[string]string)
	scanner := bufio.NewScanner(file)
	var currentSection string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = line[1 : len(line)-1]
			config[currentSection] = make(map[string]string)
			continue
		}

		if currentSection != "" {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				config[currentSection][key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return config
}

func ParseFlagTuple(tuple string) (*ParsedFlag, error) {
	tuple = strings.TrimSpace(tuple)
	if !strings.HasPrefix(tuple, "(") || !strings.HasSuffix(tuple, ")") {
		return nil, fmt.Errorf("invalid flag tuple format - must be enclosed in parentheses")
	}

	tuple = tuple[1 : len(tuple)-1]

	parts := splitTupleComponents(tuple)
	if len(parts) != 3 {
		return nil, fmt.Errorf("flag tuple must have exactly 3 components")
	}

	name := strings.TrimSpace(parts[0])
	description := strings.TrimSpace(parts[1])
	enabledStr := strings.TrimSpace(parts[2])

	if strings.HasPrefix(description, `"`) && strings.HasSuffix(description, `"`) {
		description = description[1 : len(description)-1]
	}

	enabled, err := strconv.ParseBool(enabledStr)
	if err != nil {
		return nil, fmt.Errorf("invalid boolean value in flag tuple: %v", err)
	}

	return &ParsedFlag{
		Name:        name,
		Description: description,
		Enabled:     enabled,
	}, nil
}

func splitTupleComponents(tuple string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false

	for _, r := range tuple {
		switch {
		case r == ',' && !inQuotes:
			parts = append(parts, current.String())
			current.Reset()
		case r == '"':
			inQuotes = !inQuotes
			current.WriteRune(r)
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

func GetFlagConfigs(config map[string]string) (map[string]ParsedFlag, error) {
	flags := make(map[string]ParsedFlag)

	for key, value := range config {
		if strings.HasPrefix(value, "(") && strings.HasSuffix(value, ")") {
			flag, err := ParseFlagTuple(value)
			if err != nil {
				return nil, fmt.Errorf("error parsing flag %s: %v", key, err)
			}
			flags[key] = *flag
		}
	}

	return flags, nil
}
