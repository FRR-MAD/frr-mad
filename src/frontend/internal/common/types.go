package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"net"
	"time"
)

// ================================ //
// INTERFACES                       //
// ================================ //

// PageInterface defines a module that has a title.
type PageInterface interface {
	tea.Model
	GetTitle() Tab
	GetSubTabsLength() int
	GetFooterOptions() FooterOption
}

// ================================ //
// STRUCTS                          //
// ================================ //

type WindowSize struct {
	Width  int
	Height int
}

type Tab struct {
	Title   string
	SubTabs []string
}

type FooterOption struct {
	PageTitle   string
	PageOptions []string
}

// ================================ //
// TYPES                            //
// ================================ //

// Custom type to hold and sort IPs
type IpList []string

func (ips IpList) Len() int {
	return len(ips)
}

func (ips IpList) Swap(i, j int) {
	ips[i], ips[j] = ips[j], ips[i]
}

func (ips IpList) Less(i, j int) bool {
	ip1 := net.ParseIP(ips[i])
	ip2 := net.ParseIP(ips[j])
	if ip1 == nil || ip2 == nil {
		return ips[i] < ips[j] // fallback to string comparison
	}
	return bytesCompare(ip1.To16(), ip2.To16()) < 0
}

type ReloadMessage time.Time
