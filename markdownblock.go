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
			childMd, err := child.ToMarkdown()
			if err != nil {
				return "", fmt.Errorf("failed to convert child block: %w", err)
			}
			_, err = sb.WriteString("\n\t" + childMd)
			if err != nil {
				return "", fmt.Errorf("failed to write child content: %w", err)
			}
		}
	}

	return sb.String(), nil
}
