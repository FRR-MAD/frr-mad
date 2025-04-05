package components

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func CreateTabRow(tabs []string, activeTab int, activeSubTab bool, windowSize *common.WindowSize) string {
	var renderedTabs []string
	for i, tab := range tabs {
		if i == activeTab {
			if activeSubTab {
				renderedTabs = append(renderedTabs, styles.ActiveSubTabBoxStyle.Render(tab))
			} else {
				renderedTabs = append(renderedTabs, styles.ActiveTabBoxStyle.Render(tab))
			}
		} else {
			renderedTabs = append(renderedTabs, styles.InactiveTabBoxStyle.Render(tab))
		}
	}

	// Calculate total tab row width
	tabsWidth := 0
	for _, t := range renderedTabs {
		tabsWidth += lipgloss.Width(t)
	}

	// Build the gap at the right of the last tab
	remainingWidth := max(0, windowSize.Width-tabsWidth-4)
	gap := styles.TabGap.Render(strings.Repeat(" ", remainingWidth))

	renderedTabs = append(renderedTabs, gap)

	return lipgloss.JoinHorizontal(lipgloss.Bottom, renderedTabs...)
}
