package dashboard

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var currentSubTabLocal = -1

// DashboardView is the updated View function. This allows to call View with an argument.
func (m *Model) DashboardView(currentSubTab int) string {
	currentSubTabLocal = currentSubTab
	return m.View()
}

func (m *Model) View() string {
	if currentSubTabLocal == 0 {
		return m.renderDashboardTab0()
	} else if currentSubTabLocal == 1 {
		return m.renderDashboardTab0()
	}
	return m.renderDashboardTab0()
}

func (m *Model) renderDashboardTab0() string {
	// Calculate box width dynamically for two horizontal boxes based on terminal width
	boxWidthForTwo := (m.windowSize.Width - 12) / 2 // - 6 (padding+border contentbox) - 5 (border + 1 gap)
	if boxWidthForTwo < 20 {
		boxWidthForTwo = 20 // Minimum width to ensure readability
	}

	ospfBox := styles.GeneralBoxStyle.
		Width(boxWidthForTwo).
		Render(styles.BoxTitleStyle.Render("OSPF Anomalies:") + "\n" + strings.Join(m.ospfAnomalies, "\n"))

	bgpBox := styles.BadBoxStyle.Width(boxWidthForTwo).Render(
		"Helloooo...\n\nStatus:\nYour Router explodes in 15 minutes!!!\n\nRecommendation:\nShut down everything or leave the buildings.",
	)

	horizontalBoxes := lipgloss.JoinHorizontal(lipgloss.Top, ospfBox, bgpBox)

	return horizontalBoxes
}
