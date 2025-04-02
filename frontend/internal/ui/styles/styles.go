package styles

import (
	"github.com/charmbracelet/lipgloss"
)

//var TabRowStyle = lipgloss.NewStyle().
//	Border(lipgloss.RoundedBorder()).
//	BorderBottom(false).
//	BorderForeground(lipgloss.Color("#00BFFF"))

// ----------------------------
// Box Styling
// ----------------------------

var ContentStyle = lipgloss.NewStyle().
	Border(lipgloss.DoubleBorder()).
	// BorderTop(false).
	BorderForeground(lipgloss.Color("#00BFFF")).
	Padding(0, 2)

var TitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("205")).
	Margin(1)

var BodyStyle = lipgloss.NewStyle().
	MarginLeft(2).
	MarginRight(2)

var GeneralBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#00BFFF"))

var InfoTextStyle = lipgloss.NewStyle().
	Underline(true)

var ActiveTabStyle = lipgloss.NewStyle().
	Border(lipgloss.Border{
		Top:          "─",
		Bottom:       "━",
		Left:         "│",
		Right:        "│",
		TopLeft:      "╭",
		TopRight:     "╮",
		BottomLeft:   "┗",
		BottomRight:  "┛",
		MiddleLeft:   "├",
		MiddleRight:  "┤",
		Middle:       "┼",
		MiddleTop:    "┬",
		MiddleBottom: "┴",
	}).
	BorderForeground(lipgloss.Color("#00BFFF")).
	Padding(0, 4).
	Bold(true)

var ActiveSubTabStyle = lipgloss.NewStyle().
	Border(lipgloss.Border{
		Top:          "─",
		Bottom:       "━",
		Left:         "│",
		Right:        "│",
		TopLeft:      "╭",
		TopRight:     "╮",
		BottomLeft:   "╰",
		BottomRight:  "╯",
		MiddleLeft:   "├",
		MiddleRight:  "┤",
		Middle:       "┼",
		MiddleTop:    "┬",
		MiddleBottom: "┴",
	}).
	BorderForeground(lipgloss.Color("#C0C0C0")).
	Padding(0, 4).
	Bold(true)

var InactiveTabStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#C0C0C0")).
	Padding(0, 4).
	Bold(false)
