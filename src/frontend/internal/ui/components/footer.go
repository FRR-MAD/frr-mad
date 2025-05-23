package components

import (
	"strings"
)

// Footer represents a footer component with dynamic content.
type Footer struct {
	content         []string
	mainMenuContent []string
}

// NewFooter creates a new Footer with initial content.
func NewFooter(defaultContent ...string) *Footer {
	mainMenuCopy := append([]string(nil), defaultContent...)
	return &Footer{
		content:         defaultContent,
		mainMenuContent: mainMenuCopy,
	}
}

// Append adds additional text to the footer.
func (f *Footer) Append(s string) {
	f.content = append(f.content, s)
}

// AppendMultiple adds multiple texts to the footer.
func (f *Footer) AppendMultiple(lines []string) {
	for _, s := range lines {
		f.Append(s)
	}
}

// Clean removes all entries from the footer except for the first one.
func (f *Footer) Clean() {
	if len(f.content) > 2 {
		f.content = f.content[:2]
	}
}

// Set replaces the current footer content.
func (f *Footer) Set(s string) {
	f.content = []string{s}
}

// Get returns the footer string.
func (f *Footer) Get() string {
	combined := strings.Join(f.content, " | ")
	return combined
}

// SetMainMenuOptions appends only main menu infos which were stored in defaultcontent
func (f *Footer) SetMainMenuOptions() {
	f.content = nil
	f.AppendMultiple(f.mainMenuContent)
}
