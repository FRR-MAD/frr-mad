package dashboard

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the dashboard UI.
func (m Model) View() string {
	// Use the styles defined in styles.go
	title := TitleStyle.Render(m.Title)
	metrics := strings.Join(m.Metrics, "\n")
	body := BodyStyle.Render(metrics)

	return lipgloss.JoinVertical(lipgloss.Left, title, body, "\nPress 'r' to refresh metrics.")
}
