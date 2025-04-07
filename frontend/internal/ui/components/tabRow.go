package components

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func CreateTabRow(tabs []common.Tab, activeTab int, activeSubTab int, windowSize *common.WindowSize) string {
	var renderedTabs []string
	var renderedSubTabs []string
	for i, tab := range tabs {
		if i == activeTab {
			if activeSubTab != -1 {
				renderedTabs = append(renderedTabs, styles.ActiveTabBoxLockedStyle.Render(tab.Title))
			} else {
				renderedTabs = append(renderedTabs, styles.ActiveTabBoxStyle.Render(tab.Title))
			}
			for j, subTab := range tab.SubTabs {
				if j == activeSubTab {
					renderedSubTabs = append(renderedSubTabs, styles.ActiveSubTabBoxStyle.Render(subTab))
				} else {
					renderedSubTabs = append(renderedSubTabs, styles.InactiveSubTabBoxStyle.Render(subTab))
				}
			}
		} else {
			renderedTabs = append(renderedTabs, styles.InactiveTabBoxStyle.Render(tab.Title))
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

	horizontalTabs := lipgloss.JoinHorizontal(lipgloss.Bottom, renderedTabs...)

	// subTabs := []string{"sub1", "sub2", "sub3"}

	horizontalSubTabs := lipgloss.JoinHorizontal(lipgloss.Bottom, renderedSubTabs...)

	return lipgloss.JoinVertical(lipgloss.Left, horizontalTabs, horizontalSubTabs)
}
