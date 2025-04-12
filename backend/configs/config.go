package configs

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var ConfigLocation = "/etc/frr-analytics/main.conf"

func LoadConfig() map[string]map[string]string {
	fmt.Println("Loading configuration file:", ConfigLocation)
	dat, err := os.ReadFile(ConfigLocation)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Config file %v\n", dat)

	config := make(map[string]map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(dat)))
	var title string

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// fmt.Println(line)
			title = extractConfigTitle(line, "[", "]")
			config[title] = make(map[string]string)
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		config[title][key] = value
	}

	return config
}

func extractConfigTitle(str string, start string, end string) (result string) {
	indexStart := strings.Index(str, start)
	indexEnd := strings.Index(str, end)

	return str[indexStart+1 : indexEnd]
}
