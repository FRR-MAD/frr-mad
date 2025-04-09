package configs

import (
	"fmt"
	"os"
)

var ConfigLocation = "/etc/frr-analytics/main.conf"

func LoadConfig() {
	fmt.Println("Loading configuration file:", ConfigLocation)

	dat, _ := os.ReadFile(ConfigLocation)
	fmt.Println(string(dat))

	os.Exit(0)
}
