package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// ----------------------------
// Fixed Sizes
// ----------------------------

const TabRowHeight = 6
const FooterHeight = 1

// ----------------------------
// Colors
// ----------------------------

var BoxBorderBlue = "#5f87ff" // Usage: Active Tab, Content Border
var Grey = "#444444"          // Usage: inactive components, options
var NormalBeige = "#d7d7af"   // Usage: Box Border when content good
var BadRed = "#d70000"        // Usage: Box Border when content bad
var NavyBlue = "#3a3a3a"

//var BoxBorderBlue = "111" // Usage: Active Tab, Content Border
//var Grey = "238"          // Usage: inactive components, options
//var NormalBeige = "187"   // Usage: Box Border when content good
//var BadRed = "#160"        // Usage: Box Border when content bad
//var NavyBlue = "237"

// ----------------------------
// Text Styling
// ----------------------------

var InfoTextStyle = lipgloss.NewStyle().
	Underline(true).
	Foreground(lipgloss.Color(Grey))

var BoxTitleStyle = lipgloss.NewStyle().
	Bold(true).
	Border(lipgloss.NormalBorder()).
	BorderTop(false).
	BorderLeft(false).
	BorderRight(false).
	BorderBottom(true)

var TextOutputStyle = lipgloss.NewStyle().
	Padding(1, 2)

var OSPFMonitoringTableTitleStyle = GeneralBoxStyle.
	BorderBottom(false)

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

var InactiveBoxStyle = GeneralBoxStyle.
	BorderForeground(lipgloss.Color(Grey))

var FooterBoxStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(Grey)).
	Padding(0, 1)

// ----------------------------
// Custom Borders
// ----------------------------

// If using characters for corner be aware of width deviation!

var ContentBorder = lipgloss.Border{
	Top:          " ",
	Bottom:       "─",
	Left:         "│",
	Right:        "│",
	TopLeft:      "",
	TopRight:     "",
	BottomLeft:   "╰",
	BottomRight:  "╯",
	MiddleLeft:   "├",
	MiddleRight:  "┤",
	Middle:       "┼",
	MiddleTop:    "┬",
	MiddleBottom: "┴",
}

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

var TableBorder = lipgloss.Border{
	Top:          "─",
	Bottom:       "─",
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
}

var OSPFMonitoringTableTitleBorder = lipgloss.Border{
	Top:      "─",
	Bottom:   "",
	Left:     "│",
	Right:    "│",
	TopLeft:  "╭",
	TopRight: "╮",
}

// ----------------------------
// Tab Styling
// ----------------------------

var ActiveTabBoxStyle = lipgloss.NewStyle().
	Border(ActiveTabBorder).
	BorderForeground(lipgloss.Color(BoxBorderBlue)).
	Padding(0, 4).
	Bold(true).
	Underline(true)

var ActiveTabBoxLockedStyle = ActiveTabBoxStyle.
	Bold(false).
	Underline(false)

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

var ActiveSubTabBoxStyle = lipgloss.NewStyle().
	Padding(0, 4, 0, 0).
	Bold(true).
	Underline(true)

var InactiveSubTabBoxStyle = lipgloss.NewStyle().
	Padding(0, 4, 0, 0).
	Bold(false).
	Underline(false)

//var BoxWidthForOne = lipgloss.NewStyle().
//	Width(89)

// ----------------------------
// Table Styling
// ----------------------------

var (
	HeaderStyle             = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center)
	FirstNormalRowCellStyle = lipgloss.NewStyle().Padding(0, 1)
	NormalCellStyle         = lipgloss.NewStyle().Padding(0, 1)
	BadCellStyle            = lipgloss.NewStyle().Padding(0, 1)
)

//var NormalTable = table.New().
//	Border(TableBorder).
//	BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(NormalBeige))).
//	StyleFunc(func(row, col int) lipgloss.Style {
//		switch {
//		case row == table.HeaderRow:
//			return HeaderStyle
//		default:
//			return CellStyle
//		}
//	})

//var BadTable = table.New().
//	Border(NormalTableBorder).
//	BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(BadRed))).
//	StyleFunc(func(row, col int) lipgloss.Style {
//		switch {
//		case row == table.HeaderRow:
//			return headerStyle
//		default:
//			return cellStyle
//		}
//	})
