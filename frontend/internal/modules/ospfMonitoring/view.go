package ospfMonitoring

import (
	"strings"

	"github.com/ba2025-ysmprc/frr-tui/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

// View renders the dashboard UI.
func (m Model) View() string {
	// Use the styles defined in styles.go
	title := ui.TitleStyle.Render(m.Title)
	metrics := strings.Join(m.Metrics, "\n")
	body := ui.BodyStyle.Render(metrics)

	return lipgloss.JoinVertical(lipgloss.Left, title, body, "\nPress 'r' to refresh metrics.")
}
