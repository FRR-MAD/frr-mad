package common

import tea "github.com/charmbracelet/bubbletea"

// TitledModule defines a module that has a title.
type TitledModule interface {
	tea.Model
	GetTitle() string
}
