package styles

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// ======================================== //
// Window size calculations and constants   //
// ======================================== //

const (
	ErrorPreventionContentBox = 2
	BorderContentBox          = 2
	PaddingContentBox         = 4
	BoxBorder                 = 2
	MarginX1                  = 2
	MarginX2                  = 4
	MarginX3                  = 6
	MarginX4                  = 8

	TabRowHeight = 6
	FooterHeight = 1
)

var (
	WidthBasis                    int
	roundingCorrectionOneFourthH1 int
	roundingCorrectionOneFourthH2 int

	WidthOneH1               int
	WidthOneH1Box            int
	WidthTwoH1               int
	WidthTwoH1Box            int
	WidthTwoH1OneFourth      int
	WidthTwoH1OneFourthBox   int
	WidthTwoH1ThreeFourth    int
	WidthTwoH1ThreeFourthBox int

	WidthOneH2               int
	WidthOneH2Box            int
	WidthTwoH2               int
	WidthTwoH2Box            int
	WidthTwoH2OneFourth      int
	WidthTwoH2OneFourthBox   int
	WidthTwoH2ThreeFourth    int
	WidthTwoH2ThreeFourthBox int
)

func SetWindowSizes(window common.WindowSize) {
	WidthBasis = window.Width - ErrorPreventionContentBox - BorderContentBox - PaddingContentBox
	roundingCorrectionOneFourthH1 = window.Width % 4
	roundingCorrectionOneFourthH2 = window.Width%4 - 2

	WidthOneH1 = WidthBasis - BoxBorder
	WidthOneH1Box = WidthBasis - BoxBorder - MarginX2
	WidthTwoH1 = (WidthBasis - 2*BoxBorder) / 2
	WidthTwoH1Box = (WidthBasis - 2*MarginX2) / 2
	WidthTwoH1OneFourth = (WidthBasis-2*BoxBorder)/4 + roundingCorrectionOneFourthH1
	WidthTwoH1OneFourthBox = (WidthBasis-2*MarginX2)/4 + roundingCorrectionOneFourthH1
	WidthTwoH1ThreeFourth = WidthBasis - 2*BoxBorder - WidthTwoH1OneFourth
	WidthTwoH1ThreeFourthBox = WidthBasis - 2*MarginX2 - WidthTwoH1OneFourthBox

	WidthOneH2 = WidthBasis - BoxBorder - MarginX2
	WidthOneH2Box = WidthBasis - BoxBorder - MarginX4
	WidthTwoH2 = (WidthBasis - 2*MarginX2 - 2*BoxBorder) / 2
	WidthTwoH2Box = (WidthBasis - 2*MarginX4) / 2
	WidthTwoH2OneFourth = (WidthBasis-2*MarginX2-2*BoxBorder)/4 + roundingCorrectionOneFourthH2
	WidthTwoH2OneFourthBox = (WidthBasis-2*MarginX4)/4 + roundingCorrectionOneFourthH2
	WidthTwoH2ThreeFourth = WidthBasis - 2*MarginX2 - 2*BoxBorder - WidthTwoH2OneFourth
	WidthTwoH2ThreeFourthBox = WidthBasis - 2*MarginX4 - WidthTwoH2OneFourthBox
}

// ======================================== //
// Colors                                   //
// ======================================== //

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

// ======================================== //
// Text Styling                             //
// ======================================== //

var BoxTitleStyle = lipgloss.NewStyle().
	Bold(true).
	Border(lipgloss.NormalBorder()).
	BorderTop(false).
	BorderLeft(false).
	BorderRight(false).
	BorderBottom(true)

var TextOutputStyle = lipgloss.NewStyle().
	Padding(1, 2)

func H1TitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(NormalBeige)).
		BorderBottom(false).
		Margin(0, 0, 1, 0).
		Padding(1, 0, 0, 0).
		Align(lipgloss.Center).
		Bold(true)
}

func H1TitleStyleForOne() lipgloss.Style {
	return H1TitleStyle().
		Width(WidthOneH1)
}

func H1TitleStyleForTwo() lipgloss.Style {
	return H1TitleStyleForOne().
		Width(WidthTwoH1)
}

func H2TitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(Grey)).
		BorderBottom(false).
		Margin(0, 2).
		Padding(0, 1).
		Align(lipgloss.Center).
		Bold(true)
}

func H2TitleStyleForOne() lipgloss.Style {
	return H2TitleStyle().
		Width(WidthOneH2)
}

func H2TitleStyleForTwo() lipgloss.Style {
	return H2TitleStyleForOne().
		Width(WidthTwoH2)
}

var AlignCenterAndM02P01 = lipgloss.NewStyle().
	Margin(0, 2).
	Padding(0, 1).
	Align(lipgloss.Center)

func H1BadTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(BadRed)).
		BorderBottom(false).
		Margin(0, 0, 1, 0).
		Padding(1, 0, 0, 0).
		Align(lipgloss.Center).
		Bold(true)
}

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

var GeneralBoxStyleP0M0 = GeneralBoxStyle.
	Margin(0).
	Padding(0)

var BadBoxStyle = GeneralBoxStyle.
	BorderForeground(lipgloss.Color(BadRed))

var InactiveBoxStyle = GeneralBoxStyle.
	BorderForeground(lipgloss.Color(Grey))

func H1ContentBoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Margin(0, 2).
		Padding(0, 1)
}

func H1ContentBoxCenterStyle() lipgloss.Style {
	return H1ContentBoxStyle().
		Align(lipgloss.Center)
}

func H1OneContentBoxStyle() lipgloss.Style {
	return H1ContentBoxStyle().
		Width(WidthOneH1Box)
}

func H1TwoContentBoxesStyle() lipgloss.Style {
	return H1OneContentBoxStyle().
		Width(WidthTwoH1Box)
}

func H2ContentBoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Margin(0, 4).
		Padding(0, 1)
}

func H2ContentBoxCenterStyle() lipgloss.Style {
	return H2ContentBoxStyle().
		Align(lipgloss.Center)
}

func H2OneContentBoxStyle() lipgloss.Style {
	return H2ContentBoxStyle().
		Width(WidthOneH2Box)
}

func H2OneContentBoxCenterStyle() lipgloss.Style {
	return H2OneContentBoxStyle().
		Align(lipgloss.Center)
}

func H2TwoContentBoxesStyle() lipgloss.Style {
	return H2OneContentBoxStyle().
		Width(WidthTwoH2Box)
}

func H2TwoContentBoxesCenterStyle() lipgloss.Style {
	return H2TwoContentBoxesStyle().
		Align(lipgloss.Center)
}

func H2TwoContentBoxStyleP1101() lipgloss.Style {
	return H2ContentBoxStyle().
		Padding(1, 1, 0, 1)
}

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

func H1BoxBottomBorderStyle() lipgloss.Style {
	return H1TitleStyle().
		BorderBottom(true).
		BorderTop(false)
}

func H1BadBoxBottomBorderStyle() lipgloss.Style {
	return H1BadTitleStyle().
		BorderBottom(true).
		BorderTop(false)
}

func H1OneBoxBottomBorderStyle() lipgloss.Style {
	return H1TitleStyleForOne().
		BorderBottom(true).
		BorderTop(false)
}

func H1TwoBoxBottomBorderStyle() lipgloss.Style {
	return H1TitleStyleForTwo().
		BorderBottom(true).
		BorderTop(false)
}

func H2BoxBottomBorderStyle() lipgloss.Style {
	return H2TitleStyle().
		BorderBottom(true).
		BorderTop(false)
}

func H2OneBoxBottomBorderStyle() lipgloss.Style {
	return H2TitleStyleForOne().
		BorderBottom(true).
		BorderTop(false)
}

func H2TwoBoxBottomBorderStyle() lipgloss.Style {
	return H2TitleStyleForTwo().
		BorderBottom(true).
		BorderTop(false)
}
