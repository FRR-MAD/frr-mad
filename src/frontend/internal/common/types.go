package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"net"
	"net/netip"
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

// SortedIpList holds and sorts IPs
type SortedIpList []string

func (ips SortedIpList) Len() int {
	return len(ips)
}

func (ips SortedIpList) Swap(i, j int) {
	ips[i], ips[j] = ips[j], ips[i]
}

func (ips SortedIpList) Less(i, j int) bool {
	ip1 := net.ParseIP(ips[i])
	ip2 := net.ParseIP(ips[j])
	if ip1 == nil || ip2 == nil {
		return ips[i] < ips[j] // fallback to string comparison
	}
	return bytesCompare(ip1.To16(), ip2.To16()) < 0
}

// SortedPrefixList holds and sorts IP Prefixes
type SortedPrefixList []string

func (pl SortedPrefixList) Len() int {
	return len(pl)
}

func (pl SortedPrefixList) Swap(i, j int) {
	pl[i], pl[j] = pl[j], pl[i]
}

func (pl SortedPrefixList) Less(i, j int) bool {
	p1, err1 := netip.ParsePrefix(pl[i])
	p2, err2 := netip.ParsePrefix(pl[j])

	if err1 != nil || err2 != nil {
		log.Printf("invalid prefix: %v or %v", pl[i], pl[j])
		return pl[i] < pl[j] // fallback to string comparison
	}

	// Compare IP addresses first
	if cmp := p1.Addr().Compare(p2.Addr()); cmp != 0 {
		return cmp < 0
	}

	// If addresses are equal, compare prefix length (shorter first)
	return p1.Bits() < p2.Bits()
}

type ReloadMessage time.Time
