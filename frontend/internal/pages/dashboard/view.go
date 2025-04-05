package dashboard

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

// View renders the dashboard UI.
func (m *Model) View() string {
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
		// "BGP Info:\n" + bgpInfo,
	)

	horizontalBoxes := lipgloss.JoinHorizontal(lipgloss.Top, ospfBox, bgpBox)

	//infoBox := styles.FooterBoxStyle.
	//	Width(m.windowSize.Width - 8).
	//	Render("press 'r' to refresh dashboard")
	//
	//return lipgloss.JoinVertical(lipgloss.Left, horizontalBoxes, infoBox)

	return horizontalBoxes
}
