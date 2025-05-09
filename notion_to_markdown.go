package notiontomd

import (
	"context"
	"fmt"
	"strings"

	"github.com/Kible/notion-to-md/internal/notionadapter/config"
	"github.com/Kible/notion-to-md/internal/notionadapter/gateway"
	"github.com/Kible/notion-to-md/internal/notionadapter/markdown"
	"github.com/Kible/notion-to-md/internal/notionadapter/models/requests"
	"github.com/Kible/notion-to-md/internal/notionadapter/utils"
)

type (
	Method interface {
		PageToMarkdownFull(context.Context, string) ([]*MarkdownBlock, error)
		PageToMarkdown(context.Context, string, *int) ([]*MarkdownBlock, error)
		ToMarkdownString(blocks []*MarkdownBlock) (string, error)
	}
	method struct {
		gateway *gateway.Module
		config  *Config
	}
	Params struct {
		Config *Config
	}
)

func New(p Params) (Method, error) {
	gateway, err := gateway.NewModule(&config.ConfigInternal{
		Notion: config.NotionConfig{
			Token:           p.Config.Notion.Token,
			ParseChildPages: p.Config.Notion.ParseChildPages,
		},
	})
	if err != nil {
		return nil, err
	}

	return &method{
		gateway: gateway,
		config:  p.Config,
	}, nil
}

func (m *method) PageToMarkdownFull(ctx context.Context, pageID string) ([]*MarkdownBlock, error) {
	return m.PageToMarkdown(ctx, pageID, nil)
}

func (m *method) PageToMarkdown(ctx context.Context, pageID string, pageSize *int) ([]*MarkdownBlock, error) {
	blocks, err := m.gateway.NotionAPI.GetBlockChildren(ctx, &requests.GetBlockChildrenRequest{
		BlockID:  pageID,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, err
	}

	return m.blockListToMarkdown(ctx, blocks.Results, nil, nil)
}

func (m *method) ToMarkdownString(blocks []*MarkdownBlock) (string, error) {
	if blocks == nil {
		return "", fmt.Errorf("nil markdown blocks")
	}

	var sb strings.Builder
	for _, block := range blocks {
		blockMd, err := block.ToMarkdown()
		if err != nil {
			return "", fmt.Errorf("failed to convert block: %w", err)
		}
		_, err = sb.WriteString("\n\n" + blockMd)
		if err != nil {
			return "", fmt.Errorf("failed to write block content: %w", err)
		}
	}
	return sb.String(), nil
}

func (m *method) blockListToMarkdown(ctx context.Context, blocks []any, totalPages *int, mdBlocks []*MarkdownBlock) ([]*MarkdownBlock, error) {
	if mdBlocks == nil {
		mdBlocks = make([]*MarkdownBlock, 0)
	}
	if len(blocks) == 0 {
		return mdBlocks, nil
	}

	for _, block := range blocks {
		blockMap, ok := block.(map[string]any)
		if !ok {
			continue
		}

		blockType, ok := blockMap["type"].(string)
		if !ok {
			continue
		}

		// Skip unsupported blocks and child pages if not configured to parse them
		if blockType == "unsupported" || (blockType == "child_page" && !m.config.Notion.ParseChildPages) {
			continue
		}

		// Extract block ID
		blockID := ""
		if id, ok := blockMap["id"].(string); ok {
			blockID = id
		}

		// Handle blocks with children
		if hasChildren, ok := blockMap["has_children"].(bool); ok && hasChildren {
			targetBlockID := blockID
			if blockType == "synced_block" {
				if syncedBlock, ok := blockMap["synced_block"].(map[string]any); ok {
					if syncedFrom, ok := syncedBlock["synced_from"].(map[string]any); ok {
						if id, ok := syncedFrom["block_id"].(string); ok {
							targetBlockID = id
						}
					}
				}
			}

			if targetBlockID != "" {
				childBlocks, err := m.gateway.NotionAPI.GetBlockChildren(ctx, &requests.GetBlockChildrenRequest{
					BlockID:  targetBlockID,
					PageSize: totalPages,
				})
				if err != nil {
					return nil, err
				}

				parentMarkdown, err := m.blockToMarkdown(ctx, block)
				if err != nil {
					return nil, err
				}

				mdBlock := &MarkdownBlock{
					Type:     blockType,
					BlockID:  blockID,
					Parent:   parentMarkdown,
					Children: make([]*MarkdownBlock, 0),
				}
				mdBlocks = append(mdBlocks, mdBlock)

				// Process child blocks if not a custom transformer
				if !m.isCustomTransformer(blockType) {
					children, err := m.blockListToMarkdown(ctx, childBlocks.Results, totalPages, nil)
					if err != nil {
						return nil, err
					}
					mdBlock.Children = children
				}
				continue
			}
		}

		// Handle regular blocks
		parentMarkdown, err := m.blockToMarkdown(ctx, block)
		if err != nil {
			return nil, err
		}

		mdBlock := &MarkdownBlock{
			Type:     blockType,
			BlockID:  blockID,
			Parent:   parentMarkdown,
			Children: make([]*MarkdownBlock, 0),
		}
		mdBlocks = append(mdBlocks, mdBlock)
	}

	return mdBlocks, nil
}

func (m *method) isCustomTransformer(blockType string) bool {
	// TODO: Implement custom transformer check based on configuration
	return false
}

func (m *method) blockToMarkdown(ctx context.Context, block any) (string, error) {
	blockMap, ok := block.(map[string]any)
	if !ok {
		return "", nil
	}

	blockType, ok := blockMap["type"].(string)
	if !ok {
		return "", nil
	}

	// TODO: Check for custom transformer
	// if m.isCustomTransformer(blockType) {
	//    // handle custom transformer
	// }

	parsedData := ""

	// Handle different block types
	switch blockType {
	case "image":
		if blockContent, ok := blockMap["image"].(map[string]any); ok {
			imageTitle := "image"

			// Extract caption
			captionText := ""
			if caption, ok := blockContent["caption"].([]any); ok {
				for _, item := range caption {
					if itemMap, ok := item.(map[string]any); ok {
						if plainText, ok := itemMap["plain_text"].(string); ok {
							captionText += plainText
						}
					}
				}
			}

			var link string
			imageType, _ := blockContent["type"].(string)
			if imageType == "external" {
				if external, ok := blockContent["external"].(map[string]any); ok {
					link, _ = external["url"].(string)
				}
			} else if imageType == "file" {
				if file, ok := blockContent["file"].(map[string]any); ok {
					link, _ = file["url"].(string)
				}
			}

			// Set image title from caption or URL
			if captionText != "" {
				imageTitle = strings.TrimSpace(captionText)
			} else if link != "" && strings.Contains(link, "/") {
				parts := strings.Split(link, "/")
				imageTitle = parts[len(parts)-1]
			}

			// Default to false for convertToBase64
			result, err := markdown.Image(imageTitle, link, false)
			return result, err
		}
	case "divider":
		return markdown.Divider(), nil
	case "equation":
		if equation, ok := blockMap["equation"].(map[string]any); ok {
			if expression, ok := equation["expression"].(string); ok {
				return markdown.Equation(expression), nil
			}
		}
	case "video", "file", "pdf":
		blockContent, ok := blockMap[blockType].(map[string]any)
		if ok {
			title := blockType

			// Extract caption
			captionText := ""
			if caption, ok := blockContent["caption"].([]any); ok {
				for _, item := range caption {
					if itemMap, ok := item.(map[string]any); ok {
						if plainText, ok := itemMap["plain_text"].(string); ok {
							captionText += plainText
						}
					}
				}
			}

			var link string
			fileType, _ := blockContent["type"].(string)
			if fileType == "external" {
				if external, ok := blockContent["external"].(map[string]any); ok {
					link, _ = external["url"].(string)
				}
			} else if fileType == "file" {
				if file, ok := blockContent["file"].(map[string]any); ok {
					link, _ = file["url"].(string)
				}
			}

			// Set title from caption or URL
			if captionText != "" {
				title = strings.TrimSpace(captionText)
			} else if link != "" && strings.Contains(link, "/") {
				parts := strings.Split(link, "/")
				title = parts[len(parts)-1]
			}

			return markdown.Link(title, link), nil
		}
	case "bookmark", "embed", "link_preview", "link_to_page":
		var url string
		if blockType == "link_to_page" {
			if linkToPage, ok := blockMap[blockType].(map[string]any); ok {
				if linkType, ok := linkToPage["type"].(string); ok && linkType == "page_id" {
					if pageID, ok := linkToPage["page_id"].(string); ok {
						url = fmt.Sprintf("https://www.notion.so/%s", pageID)
					}
				}
			}
		} else if blockContent, ok := blockMap[blockType].(map[string]any); ok {
			url, _ = blockContent["url"].(string)
		}

		if url != "" {
			return markdown.Link(blockType, url), nil
		}
	case "child_page":
		if !m.config.Notion.ParseChildPages {
			return "", nil
		}

		if childPage, ok := blockMap["child_page"].(map[string]any); ok {
			if pageTitle, ok := childPage["title"].(string); ok {
				// Default to returning just the heading
				return markdown.Heading2(pageTitle), nil
			}
		}
	case "child_database":
		if childDB, ok := blockMap["child_database"].(map[string]any); ok {
			pageTitle, _ := childDB["title"].(string)
			if pageTitle == "" {
				pageTitle = "child_database"
			}
			return markdown.Heading2(pageTitle), nil
		}
	case "table":
		if blockMap["has_children"] == true {
			blockID, _ := blockMap["id"].(string)
			if blockID != "" {
				tableRows := [][]string{}

				// Get table children
				tableChildren, err := m.gateway.NotionAPI.GetBlockChildren(ctx, &requests.GetBlockChildrenRequest{
					BlockID: blockID,
				})
				if err != nil {
					return "", err
				}

				// Process each row
				for _, child := range tableChildren.Results {
					if childMap, ok := child.(map[string]any); ok {
						if childType, ok := childMap["type"].(string); ok && childType == "table_row" {
							if tableRow, ok := childMap["table_row"].(map[string]any); ok {
								if cells, ok := tableRow["cells"].([]any); ok {
									row := []string{}

									// Process each cell
									for _, cell := range cells {
										// Convert cell content using BlockToMarkdown
										cellContent, err := m.blockToMarkdown(ctx, map[string]any{
											"type": "paragraph",
											"paragraph": map[string]any{
												"rich_text": cell,
											},
										})
										if err != nil {
											return "", err
										}
										row = append(row, cellContent)
									}

									tableRows = append(tableRows, row)
								}
							}
						}
					}
				}

				return markdown.Table(tableRows), nil
			}
		}
	default:
		// Rest of the types
		// "paragraph"
		// "heading_1"
		// "heading_2"
		// "heading_3"
		// "bulleted_list_item"
		// "numbered_list_item"
		// "quote"
		// "to_do"
		// "template"
		// "synced_block"
		// "child_page"
		// "child_database"
		// "code"
		// "callout"
		// "breadcrumb"
		// "table_of_contents"
		// "link_to_page"
		// "audio"
		// "unsupported"

		// Handle the rest of the block types with rich text
		blockData, ok := blockMap[blockType].(map[string]any)
		if !ok {
			return "", nil
		}

		var blockContent []any
		if text, ok := blockData["text"].([]any); ok {
			blockContent = text
		} else if richText, ok := blockData["rich_text"].([]any); ok {
			blockContent = richText
		}

		for _, content := range blockContent {
			contentMap, ok := content.(map[string]any)
			if !ok {
				continue
			}

			// Handle inline equation
			if contentType, ok := contentMap["type"].(string); ok && contentType == "equation" {
				if equation, ok := contentMap["equation"].(map[string]any); ok {
					if expression, ok := equation["expression"].(string); ok {
						parsedData += markdown.InlineEquation(expression)
					}
				}
				continue
			}

			plainText, _ := contentMap["plain_text"].(string)
			if annotations, ok := contentMap["annotations"].(map[string]any); ok {
				// Apply annotations
				if bold, ok := annotations["bold"].(bool); ok && bold {
					plainText = markdown.Bold(plainText)
				}
				if italic, ok := annotations["italic"].(bool); ok && italic {
					plainText = markdown.Italic(plainText)
				}
				if code, ok := annotations["code"].(bool); ok && code {
					plainText = markdown.InlineCode(plainText)
				}
				if strikethrough, ok := annotations["strikethrough"].(bool); ok && strikethrough {
					plainText = markdown.Strikethrough(plainText)
				}
				if underline, ok := annotations["underline"].(bool); ok && underline {
					plainText = markdown.Underline(plainText)
				}
			}

			// Add link if present
			if href, ok := contentMap["href"].(string); ok && href != "" {
				plainText = markdown.Link(plainText, href)
			}

			parsedData += plainText
		}

		// Format based on specific block types
		switch blockType {
		case "code":
			language := "text"
			if blockData, ok := blockMap["code"].(map[string]any); ok {
				if lang, ok := blockData["language"].(string); ok {
					language = lang
				}
			}
			return markdown.CodeBlock(parsedData, language), nil

		case "heading_1":
			return markdown.Heading1(parsedData), nil

		case "heading_2":
			return markdown.Heading2(parsedData), nil

		case "heading_3":
			return markdown.Heading3(parsedData), nil

		case "quote":
			return markdown.Quote(parsedData), nil

		case "callout":
			calloutString := parsedData

			// Handle callout with children
			if blockMap["has_children"] == true {
				blockID, _ := blockMap["id"].(string)
				if blockID != "" {
					// Get callout children
					calloutChildren, err := m.gateway.NotionAPI.GetBlockChildren(ctx, &requests.GetBlockChildrenRequest{
						BlockID:  blockID,
						PageSize: utils.PointerToInt(100),
					})
					if err != nil {
						return "", err
					}

					// Process callout children
					calloutChildrenMd, err := m.blockListToMarkdown(ctx, calloutChildren.Results, nil, nil)
					if err != nil {
						return "", err
					}

					calloutString += "\n"
					for _, child := range calloutChildrenMd {
						calloutString += child.Parent + "\n\n"
					}
					calloutString = strings.TrimSpace(calloutString)
				}
			}

			// Extract icon if present
			var icon *markdown.Icon
			if calloutData, ok := blockMap["callout"].(map[string]any); ok {
				if iconData, ok := calloutData["icon"].(map[string]any); ok {
					iconType, _ := iconData["type"].(string)
					icon = &markdown.Icon{
						Type: iconType,
					}

					if iconType == "emoji" {
						icon.Emoji, _ = iconData["emoji"].(string)
					} else if iconType == "external" {
						// Skip external icon data as we only need the emoji for markdown
						// icon.External = nil
					} else if iconType == "file" {
						// Skip file icon data as we only need the emoji for markdown
						// icon.File = nil
					}
				}
			}

			return markdown.Callout(calloutString, icon), nil

		case "bulleted_list_item":
			return markdown.Bullet(parsedData, nil), nil

		case "numbered_list_item":
			var number *int
			if numberData, ok := blockMap["numbered_list_item"].(map[string]any); ok {
				if n, ok := numberData["number"].(int); ok {
					number = &n
				} else if n, ok := numberData["number"].(float64); ok {
					nInt := int(n)
					number = &nInt
				}
			}
			return markdown.Bullet(parsedData, number), nil

		case "to_do":
			checked := false
			if todoData, ok := blockMap["to_do"].(map[string]any); ok {
				checked, _ = todoData["checked"].(bool)
			}
			return markdown.TODO(parsedData, checked), nil
		}
	}

	return parsedData, nil
}
