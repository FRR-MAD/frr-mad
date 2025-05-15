package toast

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

// Model represents the state of a toast notification.
type Model struct {
	text string
}

// New returns an empty toast model (no toast shown).
func New() Model {
	return Model{}
}

// showMsg signals the toast text to display.
type showMsg struct {
	Text string
}

// Show returns a command to display a toast with the given text.
func Show(text string) tea.Cmd {
	return func() tea.Msg {
		return showMsg{Text: text}
	}
}

// Update handles incoming messages for the toast model.
// On showMsg it sets the text; any other message is ignored.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case showMsg:
		m.text = msg.Text
	}
	return m, nil
}

// View renders the toast if text is non-empty; otherwise returns empty.
func (m Model) View() string {
	if m.text == "" {
		return ""
	}
	style := lipgloss.NewStyle().
		Padding(1, 1).
		Margin(1, 0).
		Background(lipgloss.Color(styles.Grey)).
		Foreground(lipgloss.Color("#ffffff"))
	return style.Render(m.text)
}

func Overlay(body string, overlay string, x, y, totalWidth, totalHeight int) string {
	// split both into lines
	bgLines := strings.Split(body, "\n")
	olLines := strings.Split(overlay, "\n")

	// make sure bgLines has exactly totalHeight lines
	for len(bgLines) < totalHeight {
		bgLines = append(bgLines, strings.Repeat(" ", totalWidth))
	}

	// splice overlay into the background lines
	for i, ol := range olLines {
		row := y + i
		if row < 0 || row >= len(bgLines) {
			continue
		}
		line := bgLines[row]
		// ensure the line is at least x columns long
		if lipgloss.Width(line) < x {
			line += strings.Repeat(" ", x-lipgloss.Width(line))
		}
		// split into before / after the overlay
		before := line[:x]
		afterStart := x + lipgloss.Width(ol)
		after := ""
		if afterStart < len(line) {
			after = line[afterStart:]
		}
		bgLines[row] = before + ol + after
	}
	return strings.Join(bgLines, "\n")
}
