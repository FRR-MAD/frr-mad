//go:build docker
// +build docker

package configs

func init() {
	ConfigLocation = "/app/config/main.yaml"
}
