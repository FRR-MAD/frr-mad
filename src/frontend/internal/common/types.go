package common

type WindowSize struct {
	Width  int
	Height int
}

type Tab struct {
	Title   string
	SubTabs []string
}

type FooterOption struct {
	PageTitle   string
	PageOptions []string
}
