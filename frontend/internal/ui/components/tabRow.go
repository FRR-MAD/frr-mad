package components

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

func CreateTabRow(tabs []string, activeTab int, activeSubTab bool) string {
	var renderedTabs []string
	for i, tab := range tabs {
		if i == activeTab {
			if activeSubTab {
				renderedTabs = append(renderedTabs, styles.ActiveSubTabStyle.Render(tab))
			} else {
				renderedTabs = append(renderedTabs, styles.ActiveTabStyle.Render(tab))
			}
		} else {
			renderedTabs = append(renderedTabs, styles.InactiveTabStyle.Render(tab))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}
