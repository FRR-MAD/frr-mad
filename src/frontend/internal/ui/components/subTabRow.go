package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/frr-mad/frr-mad/src/frontend/internal/ui/styles"
)

func CreateSubTabRow(subTabs []string, activeSubTab int) string {
	var renderedSubTabs []string
	for i, subTab := range subTabs {
		if i == activeSubTab {
			renderedSubTabs = append(renderedSubTabs, styles.ActiveSubTabBoxStyle.Render(subTab))
		} else {
			renderedSubTabs = append(renderedSubTabs, styles.InactiveSubTabBoxStyle.Render(subTab))
		}
	}

	//horizontalSubTabs := lipgloss.JoinHorizontal(lipgloss.Bottom, renderedSubTabs...)

	return lipgloss.JoinHorizontal(lipgloss.Bottom, renderedSubTabs...)
}
