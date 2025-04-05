package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// ----------------------------
// Colors
// ----------------------------

var BoxBorderBlue = "#5f87ff" // Usage: Active Tab, Content Border
var Grey = "#444444"          // Usage: inactive components, options
var NormalBeige = "#d7d7af"   // Usage: Box Border when content good
var BadRed = "#d70000"        // Usage: Box Border when content bad
var NavyBlue = "#3a3a3a"

// ----------------------------
// Text Styling
// ----------------------------

var InfoTextStyle = lipgloss.NewStyle().
	Underline(true).
	Foreground(lipgloss.Color(Grey))

var BoxTitleStyle = lipgloss.NewStyle().
	Bold(true).
	Padding(0, 0, 1, 0)

// ----------------------------
// Box Styling
// ----------------------------

var ContentBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(BoxBorderBlue)).
	Padding(0, 2)

var GeneralBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(NormalBeige)).
	Padding(0, 1)

var BadBoxStyle = GeneralBoxStyle.
	BorderForeground(lipgloss.Color(BadRed))

var FooterBoxStyle = lipgloss.NewStyle().
	//Border(FooterBoxBorder).
	//BorderForeground(lipgloss.Color(Grey)).
	Foreground(lipgloss.Color(Grey)).
	Padding(0, 1)

// ----------------------------
// Custom Borders
// ----------------------------

// If using characters for corner be aware of width deviation!

var ActiveTabBorder = lipgloss.Border{
	Top:          "─",
	Bottom:       " ",
	Left:         "│",
	Right:        "│",
	TopLeft:      "╭",
	TopRight:     "╮",
	BottomLeft:   "┘",
	BottomRight:  "└",
	MiddleLeft:   "├",
	MiddleRight:  "┤",
	Middle:       "┼",
	MiddleTop:    "┬",
	MiddleBottom: "┴",
}

var InactiveTabBorder = lipgloss.Border{
	Top:          "─",
	Bottom:       "─",
	Left:         "│",
	Right:        "│",
	TopLeft:      "╭",
	TopRight:     "╮",
	BottomLeft:   "─",
	BottomRight:  "─",
	MiddleLeft:   "├",
	MiddleRight:  "┤",
	Middle:       "┼",
	MiddleTop:    "┬",
	MiddleBottom: "┴",
}

var FooterBoxBorder = lipgloss.Border{
	Top:    "─",
	Bottom: " ",
	Left:   "",
	Right:  "",
}

// ----------------------------
// Tab Styling
// ----------------------------

var ActiveTabBoxStyle = lipgloss.NewStyle().
	Border(ActiveTabBorder).
	BorderForeground(lipgloss.Color(BoxBorderBlue)).
	Padding(0, 4).
	Bold(true)

var ActiveSubTabBoxStyle = lipgloss.NewStyle().
	Border(ActiveTabBorder).
	BorderForeground(lipgloss.Color(Grey)).
	Padding(0, 4).
	Bold(true)

var InactiveTabBoxStyle = lipgloss.NewStyle().
	Border(InactiveTabBorder).
	BorderForeground(lipgloss.Color(Grey)).
	BorderBottomForeground(lipgloss.Color(BoxBorderBlue)).
	Padding(0, 4).
	Bold(false)

var TabGap = lipgloss.NewStyle().
	Border(InactiveTabBorder).
	BorderForeground(lipgloss.Color(BoxBorderBlue)).
	BorderTop(false).
	BorderLeft(false).
	BorderRight(false)

	//var BoxWidthForOne = lipgloss.NewStyle().
	//	Width(89)
