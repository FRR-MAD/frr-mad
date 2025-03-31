package dashboard

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the dashboard UI.
func (m Model) View() string {
	// Use the styles defined in styles.go
	title := styles.TitleStyle.Render(m.Title)
	metrics := strings.Join(m.Metrics, "\n")
	body := styles.BodyStyle.Render(metrics)

	return lipgloss.JoinVertical(lipgloss.Left, title, body, "\nPress 'r' to refresh metrics.")
}
