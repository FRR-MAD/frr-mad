package components

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
	ltable "github.com/charmbracelet/lipgloss/table"
)

func NewOspfMonitorTable(headers []string, rows int) *ltable.Table {
	t := ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		BorderHeader(true).
		BorderColumn(true).
		Headers(headers...).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == ltable.HeaderRow:
				return styles.HeaderStyle
			default:
				return styles.NormalCellStyle
			}
		})
	return t
}

func NewOspfMonitorMultilineTable(headers []string, rows int) *ltable.Table {
	t := ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		BorderHeader(true).
		BorderColumn(true).
		Headers(headers...).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == ltable.HeaderRow:
				return styles.HeaderStyle
			case row == rows-1:
				return styles.LastCellOfMultiline
			default:
				return styles.MultilineCellStyle
			}
		})
	return t
}

func NewAnomalyTable(headers []string, rows int) *ltable.Table {
	t := ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		BorderHeader(true).
		BorderColumn(true).
		Headers(headers...).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == ltable.HeaderRow:
				return styles.HeaderStyle
			default:
				return styles.NormalCellStyle
			}
		})
	return t
}

func NewMultilineTable(headers []string, rows int) *ltable.Table {
	t := ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		BorderHeader(true).
		BorderColumn(true).
		Headers(headers...).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == ltable.HeaderRow:
				return styles.HeaderStyle
			case row == rows-1:
				return styles.LastCellOfMultiline
			default:
				return styles.MultilineCellStyle
			}
		})
	return t
}
