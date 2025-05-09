package models

type MarkdownBlock struct {
	Type     string
	Text     string
	Children []MarkdownBlock
}
