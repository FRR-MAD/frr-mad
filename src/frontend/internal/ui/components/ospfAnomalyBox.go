package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/frr-mad/frr-mad/src/frontend/internal/ui/styles"
)

// AnomalyBox encapsulates an OSPF anomaly box.
type AnomalyBox struct {
	Title   string
	Content string
	Style   lipgloss.Style
	Width   int
}

// NewAnomalyBox creates a new AnomalyBox with the given title, content, style, and width.
func NewAnomalyBox(title, content string, style lipgloss.Style, width int) *AnomalyBox {
	return &AnomalyBox{
		Title:   title,
		Content: content,
		Style:   style,
		Width:   width,
	}
}

// Render returns the rendered anomaly box as a string.
func (a *AnomalyBox) Render() string {
	// Combine the styled title and the content.
	boxContent := styles.TextTitleStyle.Render(a.Title) + "\n" + a.Content
	return a.Style.Width(a.Width).Render(boxContent)
}
