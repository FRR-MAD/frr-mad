package styles

import (
	"github.com/charmbracelet/bubbles/table"
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

var MainBlue = "#5f87ff"    // Usage: Active Tab, Content Border
var Grey = "#444444"        // Usage: inactive components, options
var NormalBeige = "#d7d7af" // Usage: Box Border when content good
var BadRed = "#d70000"      // Usage: Box Border when content bad
var NavyBlue = "#3a3a3a"

//var MainBlue = "111" // Usage: Active Tab, Content Border
//var Grey = "238"          // Usage: inactive components, options
//var NormalBeige = "187"   // Usage: Box Border when content good
//var BadRed = "#160"        // Usage: Box Border when content bad
//var NavyBlue = "237"

// ----------------------------
// Text Styling
// ----------------------------

var BoxTitleStyle = lipgloss.NewStyle().
	Bold(true).
	Border(lipgloss.NormalBorder()).
	BorderTop(false).
	BorderLeft(false).
	BorderRight(false).
	BorderBottom(true)

var TextOutputStyle = lipgloss.NewStyle().
	Padding(1, 2)

var H1TitleStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(NormalBeige)).
	BorderBottom(false).
	Margin(0, 0, 1, 0).
	Padding(1, 0, 0, 0).
	// Margin(0, 2).
	// Padding(0, 1).
	Align(lipgloss.Center).
	Bold(true)

var H2TitleStyle = H1TitleStyle.
	Margin(0, 2).
	Padding(0, 1).
	BorderForeground(lipgloss.Color(Grey))

var AlignCenterAndM02P01 = lipgloss.NewStyle().
	Margin(0, 2).
	Padding(0, 1).
	Align(lipgloss.Center)

// ----------------------------
// Box Styling
// ----------------------------

var ContentBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(MainBlue)).
	Padding(0, 2)

var GeneralBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(NormalBeige)).
	Margin(0, 1).
	Padding(0, 1)

var BadBoxStyle = GeneralBoxStyle.
	BorderForeground(lipgloss.Color(BadRed))

var InactiveBoxStyle = GeneralBoxStyle.
	BorderForeground(lipgloss.Color(Grey))

var H1ContentBoxStyle = lipgloss.NewStyle().
	Margin(0, 2).
	Padding(0, 1)

var H2ContentBoxStyle = lipgloss.NewStyle().
	Margin(0, 4).
	Padding(0, 1)

var H2ContentBoxStyleP1101 = H2ContentBoxStyle.
	Padding(1, 1, 0, 1)

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
	BorderForeground(lipgloss.Color(MainBlue)).
	Padding(0, 4).
	Bold(true).
	Underline(true)

var ActiveTabBoxLockedStyle = ActiveTabBoxStyle.
	Bold(false).
	Underline(false)

var InactiveTabBoxStyle = lipgloss.NewStyle().
	Border(InactiveTabBorder).
	BorderForeground(lipgloss.Color(Grey)).
	BorderBottomForeground(lipgloss.Color(MainBlue)).
	Padding(0, 4).
	Bold(false)

var TabGap = lipgloss.NewStyle().
	Border(InactiveTabBorder).
	BorderForeground(lipgloss.Color(MainBlue)).
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

func BuildTableStyles() table.Styles {
	// start from the defaults
	s := table.DefaultStyles()

	// 1) Header line: rounded border, orange text, centered, a bit of padding
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(Grey)).
		Align(lipgloss.Center).
		PaddingTop(0).
		PaddingBottom(0).
		PaddingLeft(1).
		PaddingRight(1)

	// 2) Cell style: thin border on the left, small horizontal padding
	s.Cell = s.Cell.
		BorderLeft(true).
		BorderRight(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(Grey)).
		PaddingLeft(1).
		PaddingRight(1)

	// 3) Selected row: swap fg/bg for high contrast
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(NavyBlue)).
		Bold(true)

	return s
}

// ----------------------------
// Other Styling Elements
// ----------------------------

var H1BoxBottomBorderStyle = H1TitleStyle.
	BorderBottom(true).
	BorderTop(false)

var H2BoxBottomBorderStyle = H2TitleStyle.
	BorderBottom(true).
	BorderTop(false)
