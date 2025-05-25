package common

import (
	"github.com/charmbracelet/bubbles/textinput"
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

type Filter struct {
	Query  string
	Active bool
	Input  textinput.Model
}

type FooterOption struct {
	PageTitle   string
	PageOptions []string
}

type ExportOption struct {
	Label    string
	MapKey   string
	Filename string
}

// AddExportOption adds options to ExportOption slice only if no existing entry has the same MapKey.
func AddExportOption(opts []ExportOption, opt ExportOption) []ExportOption {
	for _, e := range opts {
		if e.MapKey == opt.MapKey {
			return opts // already present
		}
	}
	return append(opts, opt)
}

type TimedPayload struct {
	ReceivedAt string      `json:"received_at"`
	Data       interface{} `json:"data"`
}

// ================================ //
// TYPES                            //
// ================================ //

// SortedIpList implements sort.Interface for a slice of IP address strings.
type SortedIpList []string

func (ips *SortedIpList) Len() int {
	return len(*ips)
}

func (ips *SortedIpList) Swap(i, j int) {
	(*ips)[i], (*ips)[j] = (*ips)[j], (*ips)[i]
}

func (ips *SortedIpList) Less(i, j int) bool {
	ip1 := net.ParseIP((*ips)[i])
	ip2 := net.ParseIP((*ips)[j])
	if ip1 == nil || ip2 == nil {
		return (*ips)[i] < (*ips)[j] // fallback to string comparison
	}
	return bytesCompare(ip1.To16(), ip2.To16()) < 0
}

// SortedPrefixList implements sort.Interface for a slice of CIDR strings.
type SortedPrefixList []string

func (pl *SortedPrefixList) Len() int {
	return len(*pl)
}

func (pl *SortedPrefixList) Swap(i, j int) {
	(*pl)[i], (*pl)[j] = (*pl)[j], (*pl)[i]
}

func (pl *SortedPrefixList) Less(i, j int) bool {
	p1, err := netip.ParsePrefix((*pl)[i])
	if err != nil {
		log.Fatalf("invalid CIDR %q: %v", (*pl)[i], err)
	}
	p2, err := netip.ParsePrefix((*pl)[j])
	if err != nil {
		log.Fatalf("invalid CIDR %q: %v", (*pl)[j], err)
	}

	if a1, a2 := p1.Addr(), p2.Addr(); a1 != a2 {
		return a1.Compare(a2) < 0
	}

	return p1.Bits() < p2.Bits()
}

type ReloadMessage time.Time

// ================================ //
// tea Messages                     //
// ================================ //

type QuitTuiFailedMsg string

func QuitTuiFailedCmd(reason string) tea.Cmd {
	return func() tea.Msg {
		return QuitTuiFailedMsg(reason)
	}
}
