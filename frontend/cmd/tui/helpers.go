package main

import (
	"os"
)

func intPtr(i int) *int {
	return &i
}

//import (
//	"fmt"
//	"github.com/ba2025-ysmprc/frr-tui/internal/ui/components"
//	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
//	// "log"
//	"os"
//
//	"github.com/ba2025-ysmprc/frr-tui/internal/pages/dashboard"
//	"github.com/ba2025-ysmprc/frr-tui/internal/pages/ospfMonitoring"
//	tea "github.com/charmbracelet/bubbletea"
//	"github.com/charmbracelet/lipgloss"
//)
//
//func (m AppModel) Titles() []string {
//	var titles []string
//	// Create a slice of TitledModule that holds all your pages.
//	pages := []tea.Model{
//		m.dashboard,
//		m.ospf,
//	}
//	for _, mod := range pages {
//		titles = append(titles, mod.Title())
//	}
//	return titles
//}

//func (m *AppModel) GetWindowSize() common.WindowSize {
//	return m.windowSize
//}

// maybeUpdateTERM updates the environment variable 'TERM' to 'xterm-256color' if necessary
func maybeUpdateTERM() {
	term := os.Getenv("TERM")
	if term == "xterm" {
		// fmt.Println("Detected TERM=xterm, updating to xterm-256color")
		os.Setenv("TERM", "xterm-256color")
	}
}
