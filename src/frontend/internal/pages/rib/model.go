package rib

import (
	"github.com/frr-mad/frr-mad/src/logger"
	"github.com/frr-mad/frr-tui/internal/common"
	backend "github.com/frr-mad/frr-tui/internal/services"
	"github.com/frr-mad/frr-tui/internal/ui/styles"
	"github.com/frr-mad/frr-tui/internal/ui/toast"
	"github.com/charmbracelet/bubbles/viewport"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	title             string
	subTabs           []string
	footer            []string
	toast             toast.Model
	cursor            int
	exportOptions     []common.ExportOption
	exportData        map[string]string
	exportDirectory   string
	showExportOverlay bool
	windowSize        *common.WindowSize
	viewport          viewport.Model
	viewportRightHalf viewport.Model
	logger            *logger.Logger
}

func New(windowSize *common.WindowSize, appLogger *logger.Logger) *Model {

	// Create the viewport with the desired dimensions.
	vp := viewport.New(styles.ViewPortWidthCompletePage, styles.ViewPortHeightCompletePage)
	vprh := viewport.New(styles.ViewPortWidthHalf, styles.ViewPortHeightCompletePage-styles.HeightH1)

	return &Model{
		title:             "RIB",
		subTabs:           []string{"RIB", "FIB", "RIB-OSPF", "RIB-BGP", "RIB-Connected", "RIB-Static"},
		footer:            []string{"[e] export options", "[r] refresh", "[↑ ↓ home end] scroll"},
		cursor:            0,
		exportOptions:     []common.ExportOption{},
		exportData:        make(map[string]string),
		exportDirectory:   "/tmp/frr-mad/exports",
		showExportOverlay: false,
		windowSize:        windowSize,
		viewport:          vp,
		viewportRightHalf: vprh,
		logger:            appLogger,
	}
}

func (m *Model) GetTitle() common.Tab {
	return common.Tab{
		Title:   m.title,
		SubTabs: m.subTabs,
	}
}

func (m *Model) GetSubTabsLength() int {
	return len(m.subTabs)
}

func (m *Model) GetFooterOptions() common.FooterOption {
	keyBoardOptions := m.footer
	return common.FooterOption{
		PageTitle:   m.title,
		PageOptions: keyBoardOptions,
	}
}

// fetchLatestData fetches all data from the backend that are possible to export from the rib exporter
func (m *Model) fetchLatestData() error {
	rib, err := backend.GetRIB(m.logger)
	if err != nil {
		return err
	}
	m.exportData["GetRIB"] = common.PrettyPrintJSON(rib)
	m.exportOptions = common.AddExportOption(m.exportOptions, common.ExportOption{
		Label:    "routing information base",
		MapKey:   "GetRIB",
		Filename: "rib.json",
	})

	ribFibSummary, err := backend.GetRibFibSummary(m.logger)
	if err != nil {
		return err
	}
	m.exportData["GetRibFibSummary"] = common.PrettyPrintJSON(ribFibSummary)
	m.exportOptions = common.AddExportOption(m.exportOptions, common.ExportOption{
		Label:    "summary data for rib and fib",
		MapKey:   "GetRibFibSummary",
		Filename: "rib_fib_summary.json",
	})

	return nil
}

func (m *Model) Init() tea.Cmd {
	return nil
}
