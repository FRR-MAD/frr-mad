package styles

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/frr-mad/frr-tui/internal/common"
)

// ======================================== //
// Window Size - calculations and constants   //
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

	TabRowHeight           = 4
	BodyFooterHeight       = 1
	FooterHeight           = 1
	AdditionalFooterHeight = 2
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

	WidthViewPortCompletePage int
	WidthViewPortHalf         int
	WidthViewPortThreeFourth  int
	WidthViewPortOneFourth    int

	HeightBasis int

	HeightViewPortCompletePage int

	HeightH1EmptyContentPadding int

	HeightH1 int
	HeightH2 int
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

	WidthViewPortCompletePage = WidthBasis + 2
	WidthViewPortHalf = WidthTwoH1 + 2
	WidthViewPortThreeFourth = WidthTwoH1ThreeFourth + 2
	WidthViewPortOneFourth = WidthTwoH1OneFourth + 2

	HeightBasis = window.Height - TabRowHeight - FooterHeight - BorderContentBox

	HeightViewPortCompletePage = HeightBasis
	HeightH1EmptyContentPadding = HeightBasis - HeightH1 - BodyFooterHeight

	HeightH1 = 4
	HeightH2 = 2
}

// ======================================== //
// Colors                                   //
// ======================================== //

var ReadModeBlue = "#5f87ff"   // Usage: Read Only Mode --> Active Menu Tab, Content Border
var WriteModeCoral = "#FF3B30" // Usage: Read/Write Mode --> Active Menu Tab, Content Border
var Grey = "#444444"           // Usage: inactive components, options, H2 Title
var NormalBeige = "#d7d7af"    // Usage: H1 Title
var GoodGreen = "#5f875f"      // Usage: Box border when content good
var BadRed = "#d70000"         // Usage: Box border when content bad
var LightBlue = "#5f87af"      // Usage: Text color to highlight every second row in a table
var NavyBlue = "#00005f"       // Usage: Text color if on NormalBeige background
var Black = "#000000"
var White = "#ffffff"

var TuiColor = ReadModeBlue

var InfoStatusColor = White
var InfoStatusBackground = Grey
var WarningStatusColor = NormalBeige
var WarningStatusBackground = Black
var ErrorStatusColor = BadRed
var ErrorStatusBackground = Black

var StatusColor = InfoStatusColor
var StatusBackground = InfoStatusBackground

func ChangeReadWriteMode(readOnlyMode bool) {
	if readOnlyMode {
		TuiColor = ReadModeBlue
	} else {
		TuiColor = WriteModeCoral
	}
}

// StatusSeverity is a simple enum for Info / Warning / Error.
type StatusSeverity int

const (
	SeverityInfo StatusSeverity = iota
	SeverityWarning
	SeverityError
)

// String implements fmt.Stringer so you can print the name if needed.
func (s StatusSeverity) String() string {
	switch s {
	case SeverityInfo:
		return "INFO"
	case SeverityWarning:
		return "WARNING"
	case SeverityError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func SetStatusSeverity(s StatusSeverity) {
	switch s {
	case SeverityInfo:
		StatusColor = InfoStatusColor
		StatusBackground = InfoStatusBackground
	case SeverityWarning:
		StatusColor = WarningStatusColor
		StatusBackground = WarningStatusBackground
	case SeverityError:
		StatusColor = ErrorStatusColor
		StatusBackground = ErrorStatusBackground
	default:
		StatusColor = InfoStatusColor
		StatusBackground = InfoStatusBackground
	}
}

// ======================================== //
// Text Styling                             //
// ======================================== //

var TextTitleStyle = lipgloss.NewStyle().
	Bold(true).
	Padding(0, 2, 0, 0).
	Border(lipgloss.NormalBorder()).
	BorderTop(false).
	BorderLeft(false).
	BorderRight(false).
	BorderBottom(true)

var TextOutputStyle = lipgloss.NewStyle().
	Padding(1, 2, 0, 2)

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

func H1GoodTitleStyle() lipgloss.Style {
	return H1TitleStyle().
		BorderForeground(lipgloss.Color(GoodGreen))
}

func H1BadTitleStyle() lipgloss.Style {
	return H1TitleStyle().
		BorderForeground(lipgloss.Color(BadRed))
}

func FilterTextStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Margin(0, 2).
		Width(WidthTwoH1Box).
		Align(lipgloss.Right)
}

func StatusTextStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		//Margin(0, 2).
		Padding(0, 1).
		Foreground(lipgloss.Color(StatusColor)).
		Background(lipgloss.Color(StatusBackground)).
		MaxWidth(WidthTwoH1Box).
		Align(lipgloss.Left)
}

var SelectedOptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(NavyBlue)).Background(lipgloss.Color(NormalBeige)).Bold(true)

// ----------------------------
// Box Styling
// ----------------------------

func ContentBoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(TuiColor)).
		Padding(0, 2)
}

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

func H1OneContentBoxCenterStyle() lipgloss.Style {
	return H1ContentBoxStyle().
		Align(lipgloss.Center).
		Width(WidthOneH1Box)
}

func H1TwoContentBoxesStyle() lipgloss.Style {
	return H1OneContentBoxStyle().
		Width(WidthTwoH1Box)
}

func H1TwoContentBoxCenterStyle() lipgloss.Style {
	return H1ContentBoxStyle().
		Align(lipgloss.Center).
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

func ActiveTabBoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(ActiveTabBorder).
		BorderForeground(lipgloss.Color(TuiColor)).
		Padding(0, 4).
		Bold(true).
		Underline(true)
}

func ActiveTabBoxLockedStyle() lipgloss.Style {
	return ActiveTabBoxStyle().
		Bold(false).
		Underline(false)
}

func InactiveTabBoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(InactiveTabBorder).
		BorderForeground(lipgloss.Color(Grey)).
		BorderBottomForeground(lipgloss.Color(TuiColor)).
		Padding(0, 4).
		Bold(false)
}

func TabGap() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(InactiveTabBorder).
		BorderForeground(lipgloss.Color(TuiColor)).
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)
}

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
	MultilineCellStyle      = lipgloss.NewStyle().Padding(0, 1, 1, 1)
	LastCellOfMultiline     = lipgloss.NewStyle().Padding(0, 1)
	BadCellStyle            = lipgloss.NewStyle().Padding(0, 1)
	EvenRowCell             = NormalCellStyle.Foreground(lipgloss.Color(LightBlue))
)

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
		Foreground(lipgloss.Color(LightBlue)).
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
func H1OneSmallBoxBottomBorderStyle() lipgloss.Style {
	return H1TitleStyleForOne().
		Padding(0).
		Margin(0).
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

var SelectedOptionCursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(Black)).Background(lipgloss.Color(NormalBeige)).Bold(true)

// ----------------------------
// Helper functions
// ----------------------------

func VerticallyCenter(content string, termHeight int) string {
	lines := lipgloss.Height(content)
	padding := (termHeight - lines) / 2
	if padding < 0 {
		padding = 0
	}
	pad := lipgloss.NewStyle().MarginTop(padding)
	return pad.Render(content)
}
