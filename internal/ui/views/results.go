package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

func DisplayResults(result map[string]interface{}) {
	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetText(fmt.Sprintf("%v", result)).
		SetDynamicColors(true).
		SetScrollable(true)

	if err := app.SetRoot(textView, true).Run(); err != nil {
		panic(err)
	}
}
