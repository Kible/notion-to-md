package markdown

import "fmt"

// InlineCode wraps text in backticks.
func InlineCode(text string) string {
	return fmt.Sprintf("`%s`", text)
}

// InlineEquation wraps text in single-dollar LaTeX delimiters.
func InlineEquation(text string) string {
	return fmt.Sprintf("$%s$", text)
}

// Bold wraps text in double-asterisks.
func Bold(text string) string {
	return fmt.Sprintf("**%s**", text)
}

// Italic wraps text in underscores.
func Italic(text string) string {
	return fmt.Sprintf("_%s_", text)
}

// Strikethrough wraps text in double-tildes.
func Strikethrough(text string) string {
	return fmt.Sprintf("~~%s~~", text)
}

// Underline uses HTML <u> tags.
func Underline(text string) string {
	return fmt.Sprintf("<u>%s</u>", text)
}

// Link creates a Markdown link.
func Link(text, href string) string {
	return fmt.Sprintf("[%s](%s)", text, href)
}
