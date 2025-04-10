package main

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	"os"
)

//func intPtr(i int) *int {
//	return &i
//}

// maybeUpdateTERM updates the environment variable 'TERM' to 'xterm-256color'
// if necessary before program start.
func maybeUpdateTERM() {
	term := os.Getenv("TERM")
	if term == "xterm" {
		// fmt.Println("Detected TERM=xterm, updating to xterm-256color")
		err := os.Setenv("TERM", "xterm-256color")
		if err != nil {
			return
		}
	}
}

// setTitles fetches all texts of all pages to fill the TabRow, SubTabRow and Footer
func (m *AppModel) setTitles() {
	pages := []common.PageInterface{
		m.dashboard,
		m.ospf,
		// code for Presentation slides
		//m.ospf2,
		//m.ospf3,
		m.shell,
	}
	for _, page := range pages {
		m.tabs = append(m.tabs, page.GetTitle())
		m.footerOptions = append(m.footerOptions, page.GetFooterOptions())
	}
}

func (m *AppModel) getCurrentFooterOptions() []string {
	for _, opt := range m.footerOptions {
		if opt.PageTitle == m.tabs[m.currentView].Title {
			return opt.PageOptions
		}
	}
	return nil
}
