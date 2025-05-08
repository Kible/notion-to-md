package markdown

import "strings"

// AddTabSpace prefixes each line with n tab characters.
func AddTabSpace(text string, n int) string {
	if n <= 0 {
		return text
	}
	tab := strings.Repeat("\t", n)
	if strings.Contains(text, "\n") {
		lines := strings.Split(text, "\n")
		for i, line := range lines {
			lines[i] = tab + line
		}
		return strings.Join(lines, "\n")
	}
	return tab + text
}
