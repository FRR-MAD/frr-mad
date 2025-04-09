package configs

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var ConfigLocation = "/etc/frr-analytics/main.conf"

func LoadConfig() map[string]string {
	fmt.Println("Loading configuration file:", ConfigLocation)
	dat, err := os.ReadFile(ConfigLocation)
	if err != nil {
		panic(err)
	}

	config := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(dat)))

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		config[key] = value
	}

	return config
}
