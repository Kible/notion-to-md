package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	notiontomd "github.com/Kible/notion-to-md"
	"github.com/joho/godotenv"
)

func TestPageToMarkdown(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			t.Fatalf("Failed to load .env file: %v", err)
		}
	}

	md, err := notiontomd.New(notiontomd.Params{
		Config: &notiontomd.Config{
			Notion: &notiontomd.NotionConfig{
				Token: os.Getenv("NOTION_API_KEY"),
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
