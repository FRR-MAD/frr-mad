package rib

import (
	"github.com/ba2025-ysmprc/frr-mad/src/logger"
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/viewport"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	title      string
	subTabs    []string
	windowSize *common.WindowSize
	viewport   viewport.Model
	logger     *logger.Logger
}

func New(windowSize *common.WindowSize, appLogger *logger.Logger) *Model {
	boxWidthForOne := windowSize.Width - 6
	// subtract tab row, footer, and border heights.
	outputHeight := windowSize.Height - styles.TabRowHeight - styles.FooterHeight - 2

	// Create the viewport with the desired dimensions.
	vp := viewport.New(boxWidthForOne, outputHeight)

	return &Model{
		title:      "RIB",
		subTabs:    []string{"RIB", "OSPF", "Connected"},
		windowSize: windowSize,
		viewport:   vp,
		logger:     appLogger,
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
	keyBoardOptions := []string{
		"[r] refresh",
		"[↑/↓] scroll",
		// "[e] export everything",
	}
	return common.FooterOption{
		PageTitle:   m.title,
		PageOptions: keyBoardOptions,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}
