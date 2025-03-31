package styles

import (
	"github.com/charmbracelet/lipgloss"
)

//var TabRowStyle = lipgloss.NewStyle().
//	Border(lipgloss.RoundedBorder()).
//	BorderBottom(false).
//	BorderForeground(lipgloss.Color("#00BFFF"))

// ContentStyle defines the style for content mounted into main frame.
var ContentStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	// BorderTop(false).
	BorderForeground(lipgloss.Color("#00BFFF")).
	Padding(0, 2)

// TitleStyle defines the style for the dashboard title.
var TitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("205")).
	Margin(1)

// BodyStyle defines the style for the metrics list.
var BodyStyle = lipgloss.NewStyle().
	MarginLeft(2).
	MarginRight(2)

var ActiveTabStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderBottom(false).
	BorderForeground(lipgloss.Color("#00BFFF")).
	Padding(0, 4).
	Bold(true)

var ActiveSubTabStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderBottom(false).
	BorderForeground(lipgloss.Color("#00BFFF")).
	Padding(0, 4).
	Bold(true)

var InactiveTabStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#C0C0C0")).
	Padding(0, 4).
	Bold(false)
