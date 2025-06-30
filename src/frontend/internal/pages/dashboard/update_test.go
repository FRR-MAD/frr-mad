package dashboard

import (
	"testing"
	// "time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestUpdate_UpKey_NoOverlay(t *testing.T) {
	// Arrange
	m := New(nil, nil, "")
	// Set known YOffset so we can see the change
	m.viewportLeft.YOffset = 5
	m.viewportRight.YOffset = 10

	// Sanity check
	if m.showAnomalyOverlay || m.showExportOverlay {
		t.Fatal("expected no overlays active")
	}

	// Act
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("up")}
	updated, _ := m.Update(msg)

	// Assert
	got := updated.(*Model)
	wantLeft := safeClamp(5 - 10)   // will be 0
	wantRight := safeClamp(10 - 10) // will be 0

	if got.viewportLeft.YOffset != wantLeft {
		t.Errorf("viewportLeft.YOffset = %d, want %d", got.viewportLeft.YOffset, wantLeft)
	}
	if got.viewportRight.YOffset != wantRight {
		t.Errorf("viewportRight.YOffset = %d, want %d", got.viewportRight.YOffset, wantRight)
	}
}

// safeClamp clamps y so it cannot be <0
func safeClamp(y int) int {
	if y < 0 {
		return 0
	}
	return y
}
