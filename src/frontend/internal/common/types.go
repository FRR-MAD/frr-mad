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

// SortedPrefixList implements sort.Interface for a slice of CIDR strings.
type SortedPrefixList []string

func (pl SortedPrefixList) Len() int {
	return len(pl)
}
func (pl SortedPrefixList) Swap(i, j int) {
	pl[i], pl[j] = pl[j], pl[i]
}

func (pl SortedPrefixList) Less(i, j int) bool {
	p1, err := netip.ParsePrefix(pl[i])
	if err != nil {
		log.Fatalf("invalid CIDR %q: %v", pl[i], err)
	}
	p2, err := netip.ParsePrefix(pl[j])
	if err != nil {
		log.Fatalf("invalid CIDR %q: %v", pl[j], err)
	}

	// 1) Compare the numeric IP addresses
	if a1, a2 := p1.Addr(), p2.Addr(); a1 != a2 {
		return a1.Compare(a2) < 0
	}

	// 2) If same IP, shorter (smaller) mask first
	return p1.Bits() < p2.Bits()
}

//func SortPrefixes(list []string) {
//	sort.Sort(SortedPrefixList(list))
//}

type ReloadMessage time.Time
