package common

import tea "github.com/charmbracelet/bubbletea"

// PageInterface defines a module that has a title.
type PageInterface interface {
	tea.Model
	GetTitle() Tab
	GetSubTabsLength() int
	GetFooterOptions() FooterOption
}
