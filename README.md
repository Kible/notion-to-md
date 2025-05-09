# Notion to Markdown Converter

A Go library for converting Notion pages and blocks to Markdown format.

## Installation

```bash
go get github.com/Kible/notion-to-md
```

## Requirements

- Go 1.24 or later
- A Notion API integration token

## Usage

### Basic Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	notiontomd "github.com/Kible/notion-to-md"
)

func main() {
	// Initialize the converter with your Notion API token
	md, err := notiontomd.New(notiontomd.Params{
		Config: &notiontomd.Config{
			Notion: &notiontomd.NotionConfig{
				Token: "your-notion-api-token"
			},
		},
	})

	if err != nil {
		log.Fatalf("Failed to create handler: %v", err)
	}

	// Convert a Notion page to Markdown
	// Use the page ID from your Notion page URL
	blocks, err := md.PageToMarkdownFull(context.Background(), "your-page-id")
	if err != nil {
		log.Fatalf("Failed to get blocks: %v", err)
	}

	mdString, err := md.ToMarkdownString(blocks)
	fmt.Println(mdString)
}
```

### Converting with Page Size Limit

```go
// If you want to limit the number of blocks fetched:
pageSize := 100 // Fetch only up to 100 blocks
blocks, err := md.PageToMarkdown(context.Background(), "your-page-id", &pageSize)
```

## API Reference

### Methods

#### `New(params notiontomd.Params) (notiontomd.Method, error)`

Creates a new instance of the Notion to Markdown converter.

#### `PageToMarkdownFull(ctx context.Context, pageID string) ([]*notiontomd.MarkdownBlock, error)`

Converts a Notion page to Markdown format, including all nested blocks.

#### `PageToMarkdown(ctx context.Context, pageID string, pageSize *int) ([]*notiontomd.MarkdownBlock, error)`

Converts a Notion page to Markdown format with an optional limit on the number of blocks.

### Types

#### `MarkdownBlock`

```go
type MarkdownBlock struct {
	Type     string           // Notion block type
	BlockID  string           // Notion block ID
	Parent   string           // Markdown content for this block
	Children []*MarkdownBlock // Child blocks (if any)
}
```

## Features

- Converts various Notion block types to Markdown
- Supports nested block structures
- Handles images, dividers, equations, videos, files, and more
- Option to parse child pages

## Getting a Notion API Token

1. Go to [https://www.notion.so/my-integrations](https://www.notion.so/my-integrations)
2. Create a new integration
3. Copy the "Internal Integration Token"
4. Share the Notion pages you want to access with your integration

## License

[license](/LICENSE)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
