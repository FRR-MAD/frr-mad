//go:build docker
// +build docker

package configs

func init() {
	ConfigLocation = "/app/local/docker-config.conf"
}
