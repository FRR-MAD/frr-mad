package ui

import "github.com/charmbracelet/lipgloss"

// TitleStyle defines the style for the dashboard title.
var TitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("205")).
	Margin(1)

// BodyStyle defines the style for the metrics list.
var BodyStyle = lipgloss.NewStyle().
	MarginLeft(2).
	MarginRight(2)
