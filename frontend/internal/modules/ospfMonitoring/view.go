package ospfMonitoring

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	title := styles.TitleStyle.Render(m.Title)
	metrics := strings.Join(m.Metrics, "\n")
	renderedMetrics := styles.BodyStyle.Render(metrics)

	// Join the title with the body, and add a prompt at the bottom.
	return lipgloss.JoinVertical(lipgloss.Left, title, renderedMetrics, "\nPress 'r' to refresh metrics.")
}
