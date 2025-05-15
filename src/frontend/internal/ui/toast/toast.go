package toast

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	text      string
	remaining int
}

// New returns an empty toast model (no toast shown).
func New() Model {
	return Model{}
}

// showMsg is an internal message to trigger showing the toast.
type showMsg struct {
	Text     string
	Duration int // seconds
}

// tickMsg is sent on each second tick to update the countdown.
type tickMsg time.Time

// Show returns a command that will display a toast with the given text
// for the specified duration. Use this in your Update to trigger a toast.
func Show(text string, duration time.Duration) tea.Cmd {
	seconds := int(duration.Seconds())
	return func() tea.Msg {
		return showMsg{Text: text, Duration: seconds}
	}
}

// Update handles toast messages: showMsg to start, tickMsg to count down.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case showMsg:
		m.text = msg.Text
		m.remaining = msg.Duration
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})

	case tickMsg:
		if m.remaining <= 0 {
			return m, nil
		}
		m.remaining--
		if m.remaining > 0 {
			return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
				return tickMsg(t)
			})
		}
		m.text = ""
		return m, nil
	}
	return m, nil
}

// View renders the toast. Returns an empty string if no toast is active.
func (m Model) View() string {
	if m.remaining <= 0 || m.text == "" {
		return ""
	}
	// style the toast box (customize as needed)
	toastStyle := lipgloss.NewStyle().
		Padding(1, 1).
		Margin(1, 0).
		Background(lipgloss.Color(styles.Grey)).
		Foreground(lipgloss.Color("#ffffff"))

	// include countdown in parentheses
	//content := fmt.Sprintf("%s (%d)", m.text, m.remaining)
	content := m.text
	return toastStyle.Render(content)
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
