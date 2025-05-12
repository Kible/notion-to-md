package notiontomd

import (
	"fmt"
	"strings"
)

type MarkdownBlock struct {
	Type     string           // Notion block type
	BlockID  string           // Notion block ID
	Parent   string           // Markdown content for this block
	Children []*MarkdownBlock // Child blocks (if any)
}

func (m *MarkdownBlock) ToMarkdown() (string, error) {
	return m.toMarkdownWithDepth(0)
}

func (m *MarkdownBlock) toMarkdownWithDepth(depth int) (string, error) {
	if m == nil {
		return "", fmt.Errorf("nil markdown block")
	}

	var sb strings.Builder
	_, err := sb.WriteString(m.Parent)
	if err != nil {
		return "", fmt.Errorf("failed to write parent content: %w", err)
	}

	if len(m.Children) > 0 {
		for _, child := range m.Children {
			childMd, err := child.toMarkdownWithDepth(depth + 1)
			if err != nil {
				return "", fmt.Errorf("failed to convert child block: %w", err)
			}
			indent := strings.Repeat("\t", depth+1)
			_, err = sb.WriteString("\n" + indent + childMd)
			if err != nil {
				return "", fmt.Errorf("failed to write child content: %w", err)
			}
		}
	}

	return sb.String(), nil
}
