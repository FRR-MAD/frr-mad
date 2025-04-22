//go:build dev
// +build dev

package configs

func init() {
	ConfigLocation = "/etc/frr_mad/frr_mad.conf"
}
