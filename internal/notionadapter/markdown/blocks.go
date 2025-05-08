package markdown

import (
	"fmt"
	"regexp"
	"strings"
)

// CodeBlock fences text with triple backticks (and optional language).
func CodeBlock(text, language string) string {
	if language == "plain text" {
		language = "text"
	}
	return fmt.Sprintf("```%s\n%s\n```", language, text)
}

// Equation wraps text in double-dollar LaTeX delimiters.
func Equation(text string) string {
	return fmt.Sprintf("$$\n%s\n$$", text)
}

// Headings level 1–3.
func Heading1(text string) string { return "# " + text }
func Heading2(text string) string { return "## " + text }
func Heading3(text string) string { return "### " + text }

// Quote prepends “> ” to each line.
func Quote(text string) string {
	spaced := strings.ReplaceAll(text, "\n", "\n> ")
	return "> " + spaced
}

// Icon represents an optional callout icon.
type Icon struct {
	Type     string            // e.g. "emoji"
	Emoji    string            // if Type=="emoji"
	External map[string]string // if Type=="external"
	File     map[string]string // if Type=="file"
}

// Callout formats a blockquote, optionally with an emoji and
// special handling if the text itself is a Markdown heading.
func Callout(text string, icon *Icon) string {
	var emoji string
	if icon != nil && icon.Type == "emoji" {
		emoji = icon.Emoji
	}

	formatted := strings.ReplaceAll(text, "\n", "\n> ")
	prefix := ""
	if emoji != "" {
		prefix = emoji + " "
	}

	// If the callout starts with a Markdown heading, preserve its level.
	re := regexp.MustCompile(`^(#{1,6})\s+(.+)`)
	if m := re.FindStringSubmatch(strings.TrimPrefix(formatted, "> ")); m != nil {
		level, content := m[1], m[2]
		hashes := strings.Repeat("#", len(level))
		return fmt.Sprintf("> %s %s%s", hashes, prefix, content)
	}

	return fmt.Sprintf("> %s%s", prefix, formatted)
}

// Bullet returns either a “- ” or “n. ” list item.
func Bullet(text string, count *int) string {
	t := strings.TrimSpace(text)
	if count != nil {
		return fmt.Sprintf("%d. %s", *count, t)
	}
	return "- " + t
}

// TODO returns a task-list item with [ ] or [x].
func TODO(text string, checked bool) string {
	mark := " "
	if checked {
		mark = "x"
	}
	return fmt.Sprintf("- [%s] %s", mark, text)
}

// Divider returns a horizontal rule.
func Divider() string {
	return "---"
}

// Toggle wraps summary/children in HTML details/summary.
func Toggle(summary, children string) string {
	if summary == "" {
		return children
	}
	return fmt.Sprintf("<details><summary>%s</summary>%s</details>", summary, children)
}
