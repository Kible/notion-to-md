package markdown

import (
	"strings"
)

// Table renders a simple Markdown table. First row is header.
func Table(cells [][]string) string {
	if len(cells) == 0 {
		return ""
	}
	headers := cells[0]
	rows := cells[1:]

	// compute column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	var b strings.Builder
	// header
	for i, h := range headers {
		b.WriteString("| ")
		b.WriteString(h)
		b.WriteString(strings.Repeat(" ", widths[i]-len(h)))
		b.WriteString(" ")
	}
	b.WriteString("|\n")
	// divider
	for _, w := range widths {
		b.WriteString("| ")
		b.WriteString(strings.Repeat("-", w))
		b.WriteString(" ")
	}
	b.WriteString("|\n")
	// rows
	for _, row := range rows {
		for i, cell := range row {
			b.WriteString("| ")
			b.WriteString(cell)
			b.WriteString(strings.Repeat(" ", widths[i]-len(cell)))
			b.WriteString(" ")
		}
		b.WriteString("|\n")
	}

	return strings.TrimSuffix(b.String(), "\n")
}
