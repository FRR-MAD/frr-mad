package main

import (
	"fmt"
	"github.com/frr-mad/frr-tui/internal/common"
	"os"
)

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

// setStartupConfig stores the running-config to check consistency before quitting FRR-MAD-TUI
func (m *AppModel) setStartupConfig() (string, error) {
	startupConfig, err := common.GetRunningConfig(m.logger)
	if err != nil {
		m.logger.Error(fmt.Sprintf("Cannot start TUI without fetching startup config. Error: %v", err))
		return "", err
	}
	return startupConfig, nil
}

// setTitles fetches all texts of all pages to fill the TabRow, SubTabRow and Footer
func (m *AppModel) setTitles() {
	pages := []common.PageInterface{
		m.dashboard,
		m.ospf,
		m.rib,
		m.shell,
	}
	for _, page := range pages {
		for _, activeView := range m.activeViews {
			if activeView == page.GetAppState() {
				m.tabs = append(m.tabs, page.GetPageInfo())
				m.footerOptions = append(m.footerOptions, page.GetFooterOptions())
			}
		}
	}
}

func (m *AppModel) getCurrentFooterOptions() []string {
	pages := []common.PageInterface{
		m.dashboard,
		m.ospf,
		m.rib,
		m.shell,
	}
	for _, opt := range m.footerOptions {
		for _, page := range pages {
			if opt.PageTitle == page.GetPageInfo().Title && m.currentView == page.GetPageInfo().AppState {
				return opt.PageOptions
			}
		}
	}
	return nil
}

func (m *AppModel) indexOfAppState(state common.AppState) int {
	for i, st := range m.activeViews {
		if st == state {
			return i
		}
	}
	return 0
}
