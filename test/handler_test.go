package test

import (
	"context"
	"fmt"
	"testing"

	notiontomd "github.com/Kible/notion-to-md"
)

func TestPageToMarkdown(t *testing.T) {
	md, err := notiontomd.New(notiontomd.Params{
		Config: &notiontomd.Config{
			Notion: &notiontomd.NotionConfig{
				Token: "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",
			},
		},
	})

	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	blocks, err := md.PageToMarkdownFull(context.Background(), "1ed57d75-c00f-808f-8c03-e803f5cc0b51")
	if err != nil {
		t.Fatalf("Failed to get blocks: %v", err)
	}

	for _, block := range blocks {
		fmt.Println(block.Parent)
	}
}
