//go:build dev
// +build dev

package configs

func init() {
	ConfigLocation = "/tmp/dev-config.yaml"
}
