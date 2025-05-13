package common

import (
	tea "github.com/charmbracelet/bubbletea"
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

type ReloadMessage time.Time
