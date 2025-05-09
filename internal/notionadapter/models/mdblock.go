package models

type MdBlock struct {
	Type     *string   `json:"type,omitempty"`
	BlockID  string    `json:"blockId"`
	Parent   string    `json:"parent"`
	Children []MdBlock `json:"children"`
}

type BlockType string

const (
	BlockImage            BlockType = "image"
	BlockVideo            BlockType = "video"
	BlockFile             BlockType = "file"
	BlockPDF              BlockType = "pdf"
	BlockTable            BlockType = "table"
	BlockBookmark         BlockType = "bookmark"
	BlockEmbed            BlockType = "embed"
	BlockEquation         BlockType = "equation"
	BlockDivider          BlockType = "divider"
	BlockToggle           BlockType = "toggle"
	BlockToDo             BlockType = "to_do"
	BlockBulletedListItem BlockType = "bulleted_list_item"
	BlockNumberedListItem BlockType = "numbered_list_item"
	BlockSyncedBlock      BlockType = "synced_block"
	BlockColumnList       BlockType = "column_list"
	BlockColumn           BlockType = "column"
	BlockLinkPreview      BlockType = "link_preview"
	BlockLinkToPage       BlockType = "link_to_page"
	BlockParagraph        BlockType = "paragraph"
	BlockHeading1         BlockType = "heading_1"
	BlockHeading2         BlockType = "heading_2"
	BlockHeading3         BlockType = "heading_3"
	BlockQuote            BlockType = "quote"
	BlockTemplate         BlockType = "template"
	BlockChildPage        BlockType = "child_page"
	BlockChildDatabase    BlockType = "child_database"
	BlockCode             BlockType = "code"
	BlockCallout          BlockType = "callout"
	BlockBreadcrumb       BlockType = "breadcrumb"
	BlockTableOfContents  BlockType = "table_of_contents"
	BlockAudio            BlockType = "audio"
	BlockUnsupported      BlockType = "unsupported"
)
